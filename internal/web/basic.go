package web

import (
	"github.com/deadblue/elevengo/internal/util"
	"io"
	"net/http"
	"strings"
)

const (
	headerContentType = "Content-Type"
)

type Payload interface {
	io.Reader
	ContentType() string
	ContentLength() int64
}

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
	return body,nil
}

func (c *Client) GetContent(url string, qs Params) (data []byte, err error) {
	body, err := c.Get(url, qs)
	if err != nil {
		return
	}
	defer util.QuietlyClose(body)
	return io.ReadAll(body)
}

func (c *Client) Touch(url string, qs Params) error {
	if body, err := c.Get(url, qs); err == nil {
		util.ConsumeReader(body)
		return nil
	} else {
		return err
	}
}

func (c *Client) Post(url string, qs Params, payload Payload) (body io.ReadCloser, err error) {
	if qs != nil {
		url = appendQueryString(url, qs)
	}
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return
	}
	req.Header.Set(headerContentType, payload.ContentType())
	if size := payload.ContentLength(); size > 0 {
		req.ContentLength = size
	}
	var resp *http.Response
	if resp, err = c.do(req); err == nil {
		body = resp.Body
	}
	return
}
