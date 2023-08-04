package elevengo

import (
	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/option"
)

// Agent holds signed-in user's credentials, and provides methods to access upstream
// server's features, such as file management, offline download, etc.
type Agent struct {
	// isWeb indicates whether the credential is for web.
	// Some API should use PC version when credential is not for web.
	isWeb bool

	// |pc| is the underlying protocol client
	pc *protocol.Client

	// Upload helper
	uh api.UploadHelper
}

// getAppVersion gets desktop client version from 115.
func (a *Agent) getAppVersion() (ver string, err error) {
	spec := (&api.AppVersionSpec{}).Init()
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	ver = spec.Result.LinuxApp.VersionCode
	return
}

// New creates Agent with customized options.
func New(options ...option.AgentOption) *Agent {
	agent := &Agent{}
	var cdMin, cdMax uint
	var name string
	// Apply options
	for _, opt := range options {
		switch opt := opt.(type) {
		case option.AgentNameOption:
			name = string(opt)
		case option.AgentCooldownOption:
			cdMin, cdMax = opt.Min, opt.Max
		case *option.AgentHttpOption:
			agent.pc = protocol.NewClient(opt.Client)
		}
	}
	if agent.pc == nil {
		agent.pc = protocol.NewClient(nil)
	}
	agent.uh.AppVer, _ = agent.getAppVersion()
	agent.pc.SetUserAgent(api.MakeUserAgent(name, agent.uh.AppVer))
	agent.pc.SetupValve(cdMin, cdMax)

	return agent
}

// Default creates an Agent with default settings.
func Default() *Agent {
	return New()
}
