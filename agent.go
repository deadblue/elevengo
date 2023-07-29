package elevengo

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/elevengo/option"
)

// Agent holds signed-in user's credentials, and provides methods to access upstream
// server's features, such as file management, offline download, etc.
type Agent struct {
	// |pc| is the underlying protocol client
	pc *protocol.Client
	// Upload helper
	uh webapi.UploadHelper
}

// getAppVersion gets desktop client version from 115.
func (a *Agent) getAppVersion() (ver string, err error) {
	qs := protocol.Params{}.
		With("callback", "get_version")
	resp := webapi.BasicResponse{}
	if err = a.pc.CallJsonpApi(webapi.ApiGetVersion, qs, &resp); err != nil {
		return
	}
	data := webapi.VersionData{}
	if err = resp.Decode(&data); err != nil {
		return
	}
	ver = data.LinuxApp.VersionCode
	return
}

// New creates Agent with customized options.
func New(options ...option.AgentOption) *Agent {
	agent := &Agent{}
	name, appVer := "", ""
	var cdMin, cdMax uint
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
	agent.pc.SetupValve(cdMin, cdMax)
	if appVer == "" {
		// Get latest app version from cloud
		appVer, _ = agent.getAppVersion()
	}
	agent.uh.SetAppVersion(appVer)
	agent.pc.SetUserAgent(webapi.MakeUserAgent(name, appVer))

	return agent
}

// Default creates an Agent with default settings.
func Default() *Agent {
	return New()
}
