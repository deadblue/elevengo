package web

import (
	"github.com/deadblue/elevengo/plugin"
	"net/http"
	"net/http/cookiejar"
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
}

func NewClient(hc plugin.HttpClient) *Client {
	client := &Client{}
	if hc == nil {
		client.cj, _ = cookiejar.New(nil)
		client.hc = defaultHttpClient(client.cj)
		client.mc = false
	} else {
		client.hc = hc
		switch hc.(type) {
		case *http.Client:
			client.cj = hc.(*http.Client).Jar
		case plugin.HttpClientWithJar:
			client.cj = hc.(plugin.HttpClientWithJar).Jar()
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
