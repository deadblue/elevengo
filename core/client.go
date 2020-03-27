package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// A production-usable HTTP client
type HttpClient interface {
	Get(url string, qs *QueryString) ([]byte, error)
}

type implHttpClient struct {
	hc         *http.Client
	beforeSend func(req *http.Request)
	afterRecv  func(resp *http.Response)
}

func (i *implHttpClient) execute(url string, qs *QueryString, form *Form) (content io.ReadCloser, err error) {
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
	method, body := http.MethodGet, io.Reader(nil)
	if form != nil {
		method, body = http.MethodPost, form.Finish()
	}
	// build request
	req, _ := http.NewRequest(method, url, body)
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
		return resp.Body, nil
	}
}

func (i *implHttpClient) Get(url string, qs *QueryString) ([]byte, error) {
	body, err := i.execute(url, qs, nil)
	if err != nil {
		return nil, err
	} else {
		return ioutil.ReadAll(body)
	}
}

func NewHttpClient(opts *HttpClientOpts) HttpClient {

	hc := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           opts.Jar,
	}

	return &implHttpClient{
		hc:         &hc,
		beforeSend: opts.BeforeSend,
		afterRecv:  opts.AfterRecv,
	}

}
