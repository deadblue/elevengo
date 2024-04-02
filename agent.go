package elevengo

import (
	"github.com/deadblue/elevengo/internal/impl"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/types"
	"github.com/deadblue/elevengo/option"
	"github.com/deadblue/elevengo/plugin"
)

// Agent holds signed-in user's credentials, and provides methods to access upstream
// server's features, such as file management, offline download, etc.
type Agent struct {
	// Low-level API client
	llc *impl.ClientImpl

	// Common parameters
	common types.CommonParams

	// isWeb indicates whether the credential is for web.
	// Some API should use PC version when credential is not for web.
	isWeb bool
}

// Default creates an Agent with default settings.
func Default() *Agent {
	return New()
}

// New creates Agent with customized options.
func New(options ...option.AgentOption) *Agent {
	var hc plugin.HttpClient = nil
	var cdMin, cdMax uint
	var name, appVer string
	// Scan options
	for _, opt := range options {
		switch opt := opt.(type) {
		case option.AgentNameOption:
			name = string(opt)
		case option.AgentCooldownOption:
			cdMin, cdMax = opt.Min, opt.Max
		case *option.AgentHttpOption:
			hc = opt.Client
		case option.AgentVersionOption:
			appVer = string(opt)
		}
	}
	llc := impl.NewClient(hc, cdMin, cdMax)
	if appVer == "" {
		appVer, _ = getLatestAppVersion(llc)
	}
	llc.SetUserAgent(protocol.MakeUserAgent(name, appVer))
	return &Agent{
		llc: llc,
		common: types.CommonParams{
			AppVer: appVer,
		},
	}
}

func getLatestAppVersion(llc client.Client) (appVer string, err error) {
	spec := (&api.AppVersionSpec{}).Init()
	if err = llc.CallApi(spec); err == nil {
		appVer = spec.Result.LinuxApp.VersionCode
	}
	return
}

// LowlevelClient returns low-level client that can directly call ApiSpec.
func (a *Agent) LowlevelClient() client.Client {
	return a.llc
}
