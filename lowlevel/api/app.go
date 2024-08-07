package api

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
)

const (
	AppAndroidLife    = "Android"
	AppAndroidTv      = "Android-tv"
	AppAndroidPan     = "115wangpan_android"
	AppBrowserWindows = "PC-115chrome"
	AppBrowserMacOS   = "MAC-115chrome"

	AppNameBrowser = "115Browser"
	AppNameDesktop = "115Desktop"
)

type AppVersionSpec struct {
	_JsonApiSpec[types.AppVersionResult, protocol.StandardResp]
}

func (s *AppVersionSpec) Init() *AppVersionSpec {
	s._JsonApiSpec.Init(
		"https://appversion.115.com/1/web/1.0/api/getMultiVer",
	)
	return s
}
