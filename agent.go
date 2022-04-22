package elevengo

import (
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/elevengo/option"
)

/*
Agent holds signed-in user's credentials, and provides methods to access upstream
server's features, such as file management, offline download, etc.
*/
type Agent struct {

	// Agent name, used info "User-Agent" header.
	name string

	// wc is the underlying web client
	wc *web.Client

	// User info
	user UserInfo

	// Offline token
	ot webapi.OfflineToken

	// Upload token
	ut webapi.UploadToken

	// hc is the old underlying http client
	// Deprecated: use wc instead.
	hc core.HttpClient
}

// New creates Agent with customized options.
func New(options ...option.Option) *Agent {
	agent := &Agent{
		name: webapi.DefaultUserAgent,
	}
	for _, opt := range options {
		switch opt.(type) {
		case option.NameOption:
			agent.name = string(opt.(option.NameOption))
		case *option.HttpOption:
			agent.wc = web.NewClient(opt.(*option.HttpOption).Client)
		}
	}
	if agent.wc == nil {
		agent.wc = web.NewClient(nil)
	}
	agent.wc.SetUserAgent(agent.name)
	return agent
}

// Default creates an Agent with default settings.
func Default() *Agent {
	return New()
}
