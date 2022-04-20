package protocol

import (
	"io"
	"net/http"
	"strings"
)

func appendQueryString(url string, qs Params) string {
	if strings.IndexRune(url, '?') < 0 {
		url = url + "?" + qs.Encode()
	} else {
		url = url + "&" + qs.Encode()
	}
	return url
}

func (c *Client) Get(url string, qs Params) (body io.ReadCloser, err error) {
	var req, resp = (*http.Request)(nil), (*http.Response)(nil)
	if qs != nil {
		url = appendQueryString(url, qs)
	}
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	}
	if resp, err = c.do(req); err == nil {
		body = resp.Body
	}
	return
}

func (c *Client) PostForm(url string, qs Params, form Params) (body io.ReadCloser, err error) {
	if qs != nil {
		url = appendQueryString(url, qs)
	}
	var reqBody io.Reader = nil
	if form != nil {
		reqBody = form.Reader()
	}
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		return
	}
	if reqBody != nil {
		req.Header.Set(headerContentType, "application/x-www-form-urlencoded")
	}
	var resp *http.Response
	if resp, err = c.do(req); err == nil {
		body = resp.Body
	}
	return
}
