package elevengo

import (
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"net/http"
)

const (
	defaultName = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:75.0) Gecko/20100101 Firefox/75.0"
	version     = "0.1.2"
)

/*
Agent holds signed-in user's credentials, and provides methods to access upstream
server's features, such as file management, offline download, etc.
*/
type Agent struct {
	name string

	hc core.HttpClient
	l  core.LoggerEx

	ui *UserInfo
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
	// additional headers
	headers := http.Header{}
	headers.Set("Accept", "*/*")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("User-Agent", name)
	return &Agent{
		name: name,
		hc:   core.NewHttpClient(headers),
		l:    core.WrapLogger(nil),
	}
}

// Create agent with default settings.
func Default() *Agent {
	return New(defaultName)
}
