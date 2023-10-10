package elevengo

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/option"
	"github.com/deadblue/elevengo/plugin"
)

// Agent holds signed-in user's credentials, and provides methods to access upstream
// server's features, such as file management, offline download, etc.
type Agent struct {
	// isWeb indicates whether the credential is for web.
	// Some API should use PC version when credential is not for web.
	isWeb bool

	// |llc| is the low-level client
	llc client.Client

	// Upload helper
	uh api.UploadHelper
}

// New creates Agent with customized options.
func New(options ...option.AgentOption) *Agent {
	var hc plugin.HttpClient = nil
	var cdMin, cdMax uint
	var name string
	// Scan options
	for _, opt := range options {
		switch opt := opt.(type) {
		case option.AgentNameOption:
			name = string(opt)
		case option.AgentCooldownOption:
			cdMin, cdMax = opt.Min, opt.Max
		case *option.AgentHttpOption:
			hc = opt.Client
		}
	}
	agent := &Agent{
		llc: protocol.NewClient(hc, cdMin, cdMax),
	}
	agent.uh.AppVer, _ = agent.getAppVersion()
	agent.llc.SetUserAgent(api.MakeUserAgent(name, agent.uh.AppVer))

	return agent
}

// Default creates an Agent with default settings.
func Default() *Agent {
	return New()
}

// getAppVersion gets desktop client version from 115.
func (a *Agent) getAppVersion() (ver string, err error) {
	spec := (&api.AppVersionSpec{}).Init()
	if err = a.llc.CallApi(spec); err != nil {
		return
	}
	ver = spec.Result.LinuxApp.VersionCode
	return
}

// LowLevelClient returns low-level client which is used by Agent.
func (a *Agent) LowLevelClient() client.Client {
	return a.llc
}
