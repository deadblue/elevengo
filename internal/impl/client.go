package impl

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/deadblue/elevengo/internal/crypto/ec115"
	"github.com/deadblue/elevengo/plugin"
)

type ClientImpl struct {
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

func NewClient(hc plugin.HttpClient, cdMin uint, cdMax uint) *ClientImpl {
	impl := &ClientImpl{
		ecc: ec115.New(),
	}
	if hc == nil {
		impl.cj, _ = cookiejar.New(nil)
		impl.hc = defaultHttpClient(impl.cj)
		impl.mc = false
	} else {
		impl.hc = hc
		switch hc := hc.(type) {
		case *http.Client:
			impl.cj = hc.Jar
		case plugin.HttpClientWithJar:
			impl.cj = hc.Jar()
		}
		if impl.cj != nil {
			impl.mc = false
		} else {
			impl.mc = true
			impl.cj, _ = cookiejar.New(nil)
		}
	}
	if cdMax > 0 && cdMax >= cdMin {
		impl.v.enabled = true
		impl.v.cdMin, impl.v.cdMax = cdMin, cdMax
	}
	return impl
}

func (c *ClientImpl) SetUserAgent(name string) {
	c.ua = name
}

func (c *ClientImpl) GetUserAgent() string {
	return c.ua
}
