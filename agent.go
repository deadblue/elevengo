package elevengo

import (
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"net/http"
	"net/http/cookiejar"
)

const (
	defaultName = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36"
	version     = "0.1.1"
)

// Agent holds an user's credentials, and provides methods to access upstream
// server's features, such as file management, offline download, etc.
type Agent struct {
	name string

	cj http.CookieJar
	hc core.HttpClient
	l  Logger

	ui *internal.UserInfo
	ot *internal.OfflineToken
}

// Get agent version.
func (a *Agent) Version() string {
	return version
}

// Create agent with specific name.
// The name will be used in User-Agent request header.
func New(name string) *Agent {
	if name == "" {
		name = defaultName
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
		name: name,
		cj:   opts.Jar,
		hc:   core.NewHttpClient(opts),
		l:    defaultLogger(),
	}
}

// Create agent with default settings.
func Default() *Agent {
	return New(defaultName)
}
