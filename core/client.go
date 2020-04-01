package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// HTTP client
type HttpClient interface {
	// Get content from a URL
	Get(url string, qs QueryString) ([]byte, error)
	// Call a JSON API with parameters
	JsonApi(url string, qs QueryString, form Form, result interface{}) (err error)
	// Call a JSON-P API with parameters
	JsonpApi(url string, qs QueryString, result interface{}) (err error)
}

type implHttpClient struct {
	hc         *http.Client
	beforeSend func(req *http.Request)
	afterRecv  func(resp *http.Response)
}

// Internal method to send request
// The returned body should be closed by invoker.
func (i *implHttpClient) send(url string, qs QueryString, form Form) (body io.ReadCloser, err error) {
	// append query string to URL
	if qs != nil {
		index := strings.IndexRune(url, '?')
		if index < 0 {
			url = fmt.Sprintf("%s?%s", url, qs.Encode())
		} else {
			url = fmt.Sprintf("%s&%s", url, qs.Encode())
		}
	}
	// process form
	method, data := http.MethodGet, io.Reader(nil)
	if form != nil {
		method, data = http.MethodPost, form.Archive()
	}
	// build request
	req, _ := http.NewRequest(method, url, data)
	if form != nil {
		req.Header.Set("Content-Type", form.ContentType())
	}
	if i.beforeSend != nil {
		i.beforeSend(req)
	}
	// send request
	if resp, err := i.hc.Do(req); err != nil {
		return nil, err
	} else {
		if i.afterRecv != nil {
			i.afterRecv(resp)
		}
		return resp.Body, nil
	}
}

func (i *implHttpClient) Get(url string, qs QueryString) ([]byte, error) {
	body, err := i.send(url, qs, nil)
	if err != nil {
		return nil, err
	}
	defer QuietlyClose(body)
	return ioutil.ReadAll(body)
}

func (i *implHttpClient) JsonApi(url string, qs QueryString, form Form, result interface{}) (err error) {
	body, err := i.send(url, qs, form)
	if err != nil {
		return
	}
	defer QuietlyClose(body)
	// parse response body
	d := json.NewDecoder(body)
	return d.Decode(result)
}

func (i *implHttpClient) JsonpApi(url string, qs QueryString, result interface{}) (err error) {
	body, err := i.send(url, qs, nil)
	if err != nil {
		return
	}
	defer QuietlyClose(body)
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	left, right := bytes.IndexByte(content, '('), bytes.LastIndexByte(content, ')')
	if left < 0 || right < 0 {
		return &json.SyntaxError{Offset: 0}
	}
	return json.Unmarshal(content[left+1:right], result)
}

type HttpOpts struct {
	// Cookie jar
	Jar http.CookieJar
	// Hook function before send a request
	BeforeSend func(req *http.Request)
	// Hook function when receive a response
	AfterRecv func(resp *http.Response)
}

func NewHttpClient(opts *HttpOpts) HttpClient {
	// Make a copy of the default tranport
	tp := http.DefaultTransport.(*http.Transport).Clone()
	// Adjust some parameters
	tp.ForceAttemptHTTP2 = true
	tp.ExpectContinueTimeout = 5 * time.Second
	tp.MaxIdleConnsPerHost = 10
	tp.MaxIdleConns = 50
	tp.IdleConnTimeout = 30 * time.Second
	tp.DialContext = (&net.Dialer{
		Timeout:   0,
		KeepAlive: 30 * time.Second,
	}).DialContext
	// Make http.Client
	hc := &http.Client{
		Transport: tp,
		Jar:       opts.Jar,
		Timeout:   0,
	}
	return &implHttpClient{
		hc:         hc,
		beforeSend: opts.BeforeSend,
		afterRecv:  opts.AfterRecv,
	}
}
