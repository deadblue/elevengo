package web

import (
	"net/http"
)

const (
	defaultUserAgent = "Mozilla/5.0"
)

func (c *Client) do(req *http.Request) (resp *http.Response, err error) {
	req.Header.Set("Accept", "*/*")
	// Always override user-agent
	ua := c.ua
	if ua == "" {
		ua = defaultUserAgent
	}
	req.Header.Set("User-Agent", ua)
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
