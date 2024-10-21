package elevengo

import (
	"context"

	"github.com/deadblue/elevengo/internal/impl"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/types"
	"github.com/deadblue/elevengo/option"
)

// Agent holds signed-in user's credentials, and provides methods to access upstream
// server's features, such as file management, offline download, etc.
type Agent struct {
	// Low-level API client
	llc *impl.ClientImpl

	// Common parameters
	common types.CommonParams
}

// Default creates an Agent with default settings.
func Default() *Agent {
	return New()
}

// New creates Agent with customized options.
func New(options ...*option.AgentOptions) *Agent {
	opts := util.NotNull(options...)
	if opts == nil {
		opts = option.Agent()
	}
	llc := impl.NewClient(opts.HttpClient, opts.CooldownMinMs, opts.CooldownMaxMs)
	appVer := opts.Version
	if appVer == "" {
		appVer, _ = getLatestAppVersion(llc, api.AppBrowserWindows)
	}
	llc.SetUserAgent(protocol.MakeUserAgent(opts.Name, api.AppNameBrowser, appVer))
	return &Agent{
		llc: llc,
		common: types.CommonParams{
			AppVer: appVer,
		},
	}
}

func getLatestAppVersion(llc client.Client, appType string) (appVer string, err error) {
	spec := (&api.AppVersionSpec{}).Init()
	if err = llc.CallApi(spec, context.Background()); err == nil {
		versionInfo := spec.Result[appType]
		appVer = versionInfo.VersionCode
	}
	return
}
