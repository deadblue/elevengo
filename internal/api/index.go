package api

import "github.com/deadblue/elevengo/internal/api/base"

type _SizeInfo struct {
	Size       float64 `json:"size"`
	SizeFormat string  `json:"size_format"`
}

type _LoginInfo struct {
	IsCurrent int    `json:"is_current"`
	LoginTime int64  `json:"utime"`
	AppFlag   string `json:"ssoent"`
	AppName   string `json:"name"`
	Ip        string `json:"ip"`
	City      string `json:"city"`
}

type IndexInfoResult struct {
	SpaceInfo struct {
		Total  _SizeInfo `json:"all_total"`
		Remain _SizeInfo `json:"all_remain"`
		Used   _SizeInfo `json:"all_use"`
	} `json:"space_info"`
	LoginInfos struct {
		List []*_LoginInfo
	} `json:"login_devices_info"`
}

type IndexInfoSpec struct {
	base.JsonApiSpec[IndexInfoResult, base.StandardResp]
}

func (s *IndexInfoSpec) Init() *IndexInfoSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/index_info")
	return s
}
