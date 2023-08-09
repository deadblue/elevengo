package protocol

import (
	"io"
	"net/http"

	"github.com/deadblue/elevengo/internal/util"
)

const (
	headerContentType = "Content-Type"
)

type Payload interface {
	io.Reader
	ContentType() string
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

func (c *Client) GetContent(url string) (data []byte, err error) {
	body, err := c.Get(url, nil)
	if err != nil {
		return
	}
	defer util.QuietlyClose(body)
	return io.ReadAll(body)
}

func (c *Client) Touch(url string) error {
	if body, err := c.Get(url, nil); err == nil {
		util.ConsumeReader(body)
		return nil
	} else {
		return err
	}
}
