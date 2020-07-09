package elevengo

import (
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"github.com/deadblue/elevengo/plugin"
	"net/http"
)

const (
	defaultName = "Mozilla/5.0"
)

/*
Agent holds signed-in user's credentials, and provides methods to access upstream
server's features, such as file management, offline download, etc.
*/
type Agent struct {
	name string

	hc core.HttpClient

	ui *UserInfo
	ot *types.OfflineToken
}

/*
Options for customize Agent.
*/
type Options struct {

	// Name of the agent, will be used in "User-Agent" request header.
	// Caller can customize it, while it does not affect any features.
	Name string

	// Logger for printing debug message.
	// Set to nil to disable the debug message.
	// Caller can implement one or simply use plugin.StdLogger.
	Logger plugin.Logger
}

// Create a customized Agent.
func New(opts *Options) *Agent {
	name, logger := defaultName, plugin.Logger(nil)
	if opts != nil {
		if len(opts.Name) > 0 {
			name = opts.Name
		}
		logger = opts.Logger
	}
	// additional headers
	headers := http.Header{}
	headers.Set("Accept", "*/*")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("User-Agent", name)
	return &Agent{
		name: name,
		hc:   core.NewHttpClient(headers, logger),
	}
}

// Create an Agent in default settings.
func Default() *Agent {
	return New(nil)
}
