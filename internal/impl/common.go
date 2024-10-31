package impl

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/client"
)

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	ne, ok := err.(net.Error)
	return ok && ne.Timeout()
}

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
	// Send request with retry
	for {
		if resp, err = c.hc.Do(req); !isTimeoutError(err) {
			break
		}
	}
	if err == nil && c.mc {
		// Save cookie
		c.cj.SetCookies(req.URL, resp.Cookies())
	}
	return
}

// |post| performs an HTTP POST request to specific URL with given payload.
func (c *ClientImpl) post(url string, payload client.Payload, context context.Context) (body io.ReadCloser, err error) {
	req, err := http.NewRequestWithContext(context, http.MethodPost, url, payload)
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

func (c *ClientImpl) Get(
	url string, headers map[string]string, context context.Context,
) (body client.Body, err error) {
	req, err := http.NewRequestWithContext(context, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	if len(headers) > 0 {
		for name, value := range headers {
			req.Header.Add(name, value)
		}
	}
	resp, err := c.send(req)
	if err == nil {
		bi := &_BodyImpl{
			rc:    resp.Body,
			size:  -1,
			total: -1,
		}
		if hv := resp.Header.Get("Content-Length"); hv != "" {
			bi.size = util.ParseInt64(hv, -1)
		}
		if hv := resp.Header.Get("Content-Range"); hv != "" {
			if index := strings.LastIndex(hv, "/"); index >= 0 {
				bi.total = util.ParseInt64(hv[index+1:], -1)
			}
		}
		body = bi
	}
	return
}
