package web

import (
	"net/http"
)

const (
	headerUserAgent = "User-Agent"
)

func (c *Client) SetUserAgent(name string) {
	c.ua = name
}

func (c *Client) do(req *http.Request) (resp *http.Response, err error) {
	// Set user agent header
	if c.ua != "" {
		req.Header.Set(headerUserAgent, c.ua)
	}
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
