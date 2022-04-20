package elevengo

import (
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/types"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/elevengo/option"
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

	pc *protocol.Client

	// hc is replaced by pc
	hc core.HttpClient

	ui *UserInfo
	ot *types.OfflineToken
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
			agent.pc = protocol.NewClient(opt.(*option.HttpOption).Client)
		}
	}
	if agent.pc == nil {
		agent.pc = protocol.NewClient(nil)
	}
	agent.pc.SetUserAgent(agent.name)
	return agent
}

// Default creates an Agent with default settings.
func Default() *Agent {
	return New()
}
