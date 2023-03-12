package web

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
