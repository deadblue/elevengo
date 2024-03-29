package types

type SizeInfo struct {
	Size       float64 `json:"size"`
	SizeFormat string  `json:"size_format"`
}

type LoginInfo struct {
	IsCurrent int    `json:"is_current"`
	LoginTime int64  `json:"utime"`
	AppFlag   string `json:"ssoent"`
	AppName   string `json:"name"`
	Ip        string `json:"ip"`
	City      string `json:"city"`
}

type IndexInfoResult struct {
	SpaceInfo struct {
		Total  SizeInfo `json:"all_total"`
		Remain SizeInfo `json:"all_remain"`
		Used   SizeInfo `json:"all_use"`
	} `json:"space_info"`
	LoginInfos struct {
		List []*LoginInfo
	} `json:"login_devices_info"`
}
