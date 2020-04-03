package elevengo

import (
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"net/http"
	"net/http/cookiejar"
)

const (
	Version = "0.1.1"

	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"
)

type Agent struct {
	cj http.CookieJar
	hc core.HttpClient
	ua string

	ui *internal.UserInfo
	ot *internal.OfflineToken
}

func New(name string) *Agent {
	if name == "" {
		name = defaultUserAgent
	}
	opts := &core.HttpOpts{}
	opts.Jar, _ = cookiejar.New(nil)
	opts.BeforeSend = func(req *http.Request) {
		// Set headers
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("User-Agent", name)
	}
	return &Agent{
		ua: name,
		cj: opts.Jar,
		hc: core.NewHttpClient(opts),
	}
}

func Default() *Agent {
	return New(defaultUserAgent)
}
