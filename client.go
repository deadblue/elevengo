package elevengo

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type _UserInfo struct {
	UserId string
}

type _OfflineToken struct {
	Sign string
	Time int64
}

type Options struct {
	UserAgent    string
	DisableProxy bool
	MaxIdleConns int
	Debug        bool
}

type Client struct {
	jar    http.CookieJar
	client *http.Client

	userAgent string
	info      *_UserInfo
	offline   *_OfflineToken
}

func New(opts *Options) *Client {
	// core component
	dialer := &net.Dialer{
		Timeout: 30 * time.Second,
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   100,
		DialContext:           dialer.DialContext,
		IdleConnTimeout:       300 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	jar, _ := cookiejar.New(nil)
	hc := &http.Client{
		Transport: transport,
		Jar:       jar,
	}
	// assemble the client
	client := &Client{
		jar:       jar,
		client:    hc,
		userAgent: defaultUserAgent,
	}
	// apply options
	if opts != nil {
		if opts.UserAgent != "" {
			client.userAgent = opts.UserAgent
		}
		if opts.DisableProxy {
			transport.Proxy = func(request *http.Request) (url *url.URL, e error) {
				return nil, nil
			}
		}
		if opts.MaxIdleConns > 0 {
			transport.MaxConnsPerHost = opts.MaxIdleConns
			transport.MaxIdleConns = 2 * opts.MaxIdleConns
		}
		if opts.Debug {
			transport.TLSClientConfig.InsecureSkipVerify = true
		}
	}
	return client
}

func Default() *Client {
	return New(nil)
}
