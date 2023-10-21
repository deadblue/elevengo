package impl

import (
	"io"
	"net/http"

	"github.com/deadblue/elevengo/lowlevel/client"
)

// |send| sends an HTTP request, returns HTTP response or an error.
func (c *ClientImpl) send(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("Accept", "*/*")
	// Always override user-agent
	ua := c.ua
	if ua == "" {
		ua = defaultUserAgent
	}
	req.Header.Set(headerUserAgent, ua)
	if c.mc {
		// Add cookie
		for _, cookie := range c.cj.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}
	if resp, err = c.hc.Do(req); err == nil {
		if c.mc {
			// Save cookie
			c.cj.SetCookies(req.URL, resp.Cookies())
		}
	}
	return
}

// |post| performs an HTTP POST request to specific URL with given payload.
func (c *ClientImpl) post(url string, payload client.Payload) (body io.ReadCloser, err error) {
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return
	}
	req.Header.Set(headerContentType, payload.ContentType())
	if size := payload.ContentLength(); size > 0 {
		req.ContentLength = size
	}
	var resp *http.Response
	if resp, err = c.send(req); err == nil {
		body = resp.Body
	}
	return
}

func (c *ClientImpl) Get(url string, headers map[string]string) (body io.ReadCloser, err error) {
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
	if resp, err = c.send(req); err == nil {
		body = resp.Body
	}
	return
}
