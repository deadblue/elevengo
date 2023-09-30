package protocol

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/deadblue/elevengo/internal/crypto/ec115"
	"github.com/deadblue/elevengo/plugin"
)

type Client struct {
	// HTTP client
	hc plugin.HttpClient
	// Cookie jar
	cj http.CookieJar
	// Should Client manage cookie
	mc bool

	// User agent
	ua string

	// EC cipher
	ecc *ec115.Cipher

	// Valve to control API request frequency
	v Valve
}

func NewClient(hc plugin.HttpClient) *Client {
	client := &Client{
		ecc: ec115.New(),
	}
	if hc == nil {
		client.cj, _ = cookiejar.New(nil)
		client.hc = defaultHttpClient(client.cj)
		client.mc = false
	} else {
		client.hc = hc
		switch hc := hc.(type) {
		case *http.Client:
			client.cj = hc.Jar
		case plugin.HttpClientWithJar:
			client.cj = hc.Jar()
		}
		if client.cj != nil {
			client.mc = false
		} else {
			client.mc = true
			client.cj, _ = cookiejar.New(nil)
		}
	}
	return client
}

func (c *Client) SetUserAgent(name string) {
	c.ua = name
}

func (c *Client) GetUserAgent() string {
	return c.ua
}

func (c *Client) SetupValve(cdMin, cdMax uint) {
	if cdMax > 0 && cdMax >= cdMin {
		c.v.enabled = true
		c.v.cdMin, c.v.cdMax = cdMin, cdMax
	}
}
