package types

type AppVersionInfo struct {
	AppOs       int    `json:"app_os"`
	CreatedTime int64  `json:"created_time"`
	VersionCode string `json:"version_code"`
	VersionUrl  string `json:"version_url"`
}

type AppVersionResult map[string]*AppVersionInfo
