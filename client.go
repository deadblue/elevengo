package elevengo

import (
	"net"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	jar http.CookieJar
	hc  *http.Client
	ua  string

	info    *_UserInfo
	offline *_OfflineToken
}

func New(opts *Options) *Client {
	if opts == nil {
		opts = NewOptions()
	}
	// core component
	d := &net.Dialer{
		Timeout: opts.DialTimeout,
	}
	tp := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         d.DialContext,
		IdleConnTimeout:     opts.IdleTimeout,
		MaxIdleConnsPerHost: opts.MaxIdleConnsPreHost,
		MaxIdleConns:        opts.MaxIdleConns,
	}
	jar, _ := cookiejar.New(nil)
	hc := &http.Client{
		Transport: tp,
		Jar:       jar,
	}
	// assemble the client
	return &Client{
		jar: jar,
		hc:  hc,
		ua:  opts.UserAgent,
	}
}

func Default() *Client {
	return New(nil)
}
