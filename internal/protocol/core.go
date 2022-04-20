package protocol

import (
	"net/http"
)

const (
	headerContentType = "Content-Type"

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
	// Add cookie
	for _, cookie := range c.cj.Cookies(req.URL) {
		req.AddCookie(cookie)
	}
	return c.hc.Do(req)
}
