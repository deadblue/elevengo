package elevengo

import (
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"net/http"
	"net/http/cookiejar"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"
)

type Client struct {
	cj http.CookieJar
	hc core.HttpClient
	ua string

	ui *internal.UserInfo
	ot *internal.OfflineToken
}

func New(userAgent string) *Client {
	if userAgent == "" {
		userAgent = defaultUserAgent
	}
	opts := core.NewHttpOpts()
	opts.Jar, _ = cookiejar.New(nil)
	opts.BeforeSend = func(req *http.Request) {
		// Set headers
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("User-Agent", userAgent)
	}
	return &Client{
		ua: userAgent,
		cj: opts.Jar,
		hc: core.NewHttpClient(opts),
	}
}

func Default() *Client {
	return New("")
}
