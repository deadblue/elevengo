package webapi

type VersionInfo struct {
	CreatedTime int64  `json:"created_time"`
	VersionCode string `json:"version_code"`
	VersionUrl  string `json:"version_url"`
}

type VersionData struct {
	Android    VersionInfo `json:"android"`
	LinuxApp   VersionInfo `json:"linux_115"`
	MacBrowser VersionInfo `json:"mac"`
	MacApp     VersionInfo `json:"mac_115"`
	WinBrowser VersionInfo `json:"win"`
	WinApp     VersionInfo `json:"window_115"`
}
