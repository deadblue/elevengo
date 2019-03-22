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

type Client struct {
	jar http.CookieJar
	hc  *http.Client
	ua  string

	info    *_UserInfo
	offline *_OfflineToken
	captcha *_CaptchaSession
}

type Options struct {
	// The UserAgent
	UserAgent string

	// Do not use proxy in environment
	DisableProxy bool

	// Max idle connections number per host
	MaxIdleConns int

	// Enable debug mode
	// When enabled, the client:
	//  * Does not verify server certificate
	Debug bool
}

func New(opts *Options) *Client {
	// core component
	dialer := &net.Dialer{
		Timeout: defaultConnTimeout,
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   defaultConnTimeout,
		ResponseHeaderTimeout: defaultServerTimeout,
		IdleConnTimeout:       defaultIdleTimeout,
		MaxIdleConnsPerHost:   defaultIdleConnsPreHost,
		MaxIdleConns:          200,
		ExpectContinueTimeout: 1 * time.Second,
	}
	jar, _ := cookiejar.New(nil)
	hc := &http.Client{
		Transport: transport,
		Jar:       jar,
	}
	// assemble the client
	client := &Client{
		jar: jar,
		hc:  hc,
		ua:  defaultUserAgent,
	}
	// apply options
	if opts != nil {
		if opts.UserAgent != "" {
			client.ua = opts.UserAgent
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
