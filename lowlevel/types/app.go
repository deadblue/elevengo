package types

type AppVersionInfo struct {
	CreatedTime int64  `json:"created_time"`
	VersionCode string `json:"version_code"`
	VersionUrl  string `json:"version_url"`
}

type AppVersionResult struct {
	Android    AppVersionInfo `json:"android"`
	LinuxApp   AppVersionInfo `json:"linux_115"`
	MacBrowser AppVersionInfo `json:"mac"`
	MacApp     AppVersionInfo `json:"mac_115"`
	WinBrowser AppVersionInfo `json:"win"`
	WinApp     AppVersionInfo `json:"window_115"`
}
