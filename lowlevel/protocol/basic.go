package protocol

import (
	"io"
	"net/http"
)

const (
	headerContentType = "Content-Type"
)

// Payload describes the request body.
type Payload interface {
	io.Reader
	// ContentType returns the MIME type of payload.
	ContentType() string
	// ContentLength returns the size in bytes of payload.
	ContentLength() int64
}

func (c *Client) Get(url string, headers map[string]string) (body io.ReadCloser, err error) {
	var req *http.Request = nil
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		return
	}
	if len(headers) > 0 {
		for name, value := range headers {
			req.Header.Add(name, value)
		}
	}
	var resp *http.Response = nil
	if resp, err = c.do(req); err == nil {
		body = resp.Body
	}
	return
}

func (c *Client) Post(url string, payload Payload) (body io.ReadCloser, err error) {
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
