package api

import "github.com/deadblue/elevengo/internal/api/base"

type _AppVersionInfo struct {
	CreatedTime int64  `json:"created_time"`
	VersionCode string `json:"version_code"`
	VersionUrl  string `json:"version_url"`
}

type AppVersionResult struct {
	Android    _AppVersionInfo `json:"android"`
	LinuxApp   _AppVersionInfo `json:"linux_115"`
	MacBrowser _AppVersionInfo `json:"mac"`
	MacApp     _AppVersionInfo `json:"mac_115"`
	WinBrowser _AppVersionInfo `json:"win"`
	WinApp     _AppVersionInfo `json:"window_115"`
}

type AppVersionSpec struct {
	base.JsonpApiSpec[AppVersionResult, base.StandardResp]
}

func (s *AppVersionSpec) Init() *AppVersionSpec {
	s.JsonpApiSpec.Init(
		"https://appversion.115.com/1/web/1.0/api/chrome", "get_version",
	)
	return s
}
