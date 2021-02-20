package mobile

import (
	"github.com/deadblue/elevengo/internal/cipher/ec115"
	"net/http"
)

type Client struct {
	// HTTP client
	hc *http.Client
	// Cookie jar
	cj http.CookieJar
	// User ID
	uid uint32

	// EC coder
	ec *ec115.Coder
}

func New() (client *Client, err error) {
	client = &Client{
		ec: ec115.New(),
	}
	client.httpInit()

	return
}
