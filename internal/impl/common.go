package impl

import (
	"context"
	"net"
	"net/http"

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

func (c *ClientImpl) Post(
	url string, payload client.Payload, headers map[string]string, context context.Context,
) (body client.Body, err error) {
	req, err := http.NewRequestWithContext(context, http.MethodPost, url, payload)
	if err != nil {
		return
	}
	if len(headers) > 0 {
		for name, value := range headers {
			req.Header.Add(name, value)
		}
	}
	req.Header.Set(headerContentType, payload.ContentType())
	if size := payload.ContentLength(); size > 0 {
		req.ContentLength = size
	}
	var resp *http.Response
	if resp, err = c.send(req); err == nil {
		body = makeClientBody(resp)
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
	if resp, err := c.send(req); err == nil {
		body = makeClientBody(resp)
	}
	return
}

func makeClientBody(resp *http.Response) client.Body {
	body := &_BodyImpl{
		rc:   resp.Body,
		size: -1,
	}
	if hv := resp.Header.Get("Content-Length"); hv != "" {
		body.size = util.ParseInt64(hv, -1)
	}
	return body
}
