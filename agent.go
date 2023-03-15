package elevengo

import (
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/elevengo/option"
)

// Agent holds signed-in user's credentials, and provides methods to access upstream
// server's features, such as file management, offline download, etc.
type Agent struct {
	// wc is the underlying web client
	wc *web.Client

	// User ID
	uid int
	// Offline token
	ot webapi.OfflineToken
	// Upload helper
	uh webapi.UploadHelper
}

// getAppVersion gets desktop client version from 115.
func (a *Agent) getAppVersion() (ver string, err error) {
	qs := web.Params{}.
		With("callback", "get_version")
	resp := webapi.BasicResponse{}
	if err = a.wc.CallJsonpApi(webapi.ApiGetVersion, qs, &resp); err != nil {
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
func New(options ...option.Option) *Agent {
	agent, name := &Agent{}, ""
	// Apply options
	for _, opt := range options {
		switch opt := opt.(type) {
		case option.NameOption:
			name = string(opt)
		case *option.HttpOption:
			agent.wc = web.NewClient(opt.Client)
		}
	}
	if agent.wc == nil {
		agent.wc = web.NewClient(nil)
	}

	// Get latest app version from cloud
	appVer, _ := agent.getAppVersion()
	agent.uh.SetAppVersion(appVer)
	agent.wc.SetUserAgent(webapi.MakeUserAgent(name, appVer))
	
	return agent
}

// Default creates an Agent with default settings.
func Default() *Agent {
	return New()
}
