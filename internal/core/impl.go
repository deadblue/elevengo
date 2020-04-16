package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deadblue/elevengo/internal/util"
	"io"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"strings"
)

type implHttpClient struct {
	// Underlying HTTP transport client.
	hc *http.Client
	// Cookie storage.
	jar http.CookieJar
	// Additional request headers that will be sent in every request.
	hdrs http.Header
	// logger.
	l LoggerEx
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
	i.l.Printf("Request %s => %s", method, url)
	// build request
	req, _ := http.NewRequest(method, url, data)
	if form != nil {
		req.Header.Set("Content-Type", form.ContentType())
	}
	// Add additional request headers
	if i.hdrs != nil {
		for key := range i.hdrs {
			req.Header.Set(key, i.hdrs.Get(key))
		}
	}
	// send request
	if resp, err := i.hc.Do(req); err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

func (i *implHttpClient) Get(url string, qs QueryString) ([]byte, error) {
	body, err := i.send(url, qs, nil)
	if err != nil {
		return nil, err
	}
	defer util.QuietlyClose(body)
	return ioutil.ReadAll(body)
}

func (i *implHttpClient) JsonApi(url string, qs QueryString, form Form, result interface{}) (err error) {
	body, err := i.send(url, qs, form)
	if err != nil {
		return
	}
	defer util.QuietlyClose(body)
	// parse response body
	d := json.NewDecoder(body)
	return d.Decode(result)
}

func (i *implHttpClient) JsonpApi(url string, qs QueryString, result interface{}) (err error) {
	body, err := i.send(url, qs, nil)
	if err != nil {
		return
	}
	defer util.QuietlyClose(body)
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

func (i *implHttpClient) SetCookies(url, domain string, cookies map[string]string) {
	if u, err := neturl.Parse(url); err != nil {
		return
	} else {
		cks, index := make([]*http.Cookie, len(cookies)), 0
		for name, value := range cookies {
			cks[index] = &http.Cookie{
				Name:     name,
				Value:    value,
				Domain:   domain,
				Path:     "/",
				HttpOnly: true,
			}
			index += 1
		}
		i.jar.SetCookies(u, cks)
	}
}

func (i *implHttpClient) Cookies(url string) (cookies map[string]string) {
	cookies = make(map[string]string)
	if u, err := neturl.Parse(url); err != nil {
		return
	} else {
		for _, ck := range i.jar.Cookies(u) {
			cookies[ck.Name] = ck.Value
		}
	}
	return
}
