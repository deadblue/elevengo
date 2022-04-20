package protocol

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
	// User agent
	ua string
}

func NewClient(hc plugin.HttpClient) *Client {
	client := &Client{}
	if hc == nil {
		client.cj, _ = cookiejar.New(nil)
		client.hc = defaultHttpClient(client.cj)
	} else {
		switch hc.(type) {
		case *http.Client:
			client.cj = hc.(*http.Client).Jar
		case plugin.HttpClientWithJar:
			client.cj = hc.(plugin.HttpClientWithJar).Jar()
		}
		if client.cj == nil {
			client.cj, _ = cookiejar.New(nil)
		}
	}
	return client
}
