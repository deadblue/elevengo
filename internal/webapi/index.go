package webapi

type LoginDeviceInfo struct {
	SsoId     string `json:"ssoent"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Device    string `json:"device"`
	Ip        string `json:"ip"`
	City      string `json:"city"`
	Time      int64  `json:"utime"`
	IsCurrent int    `json:"is_current"`
}

type SizeInfo struct {
	Size       float64 `json:"size"`
	FormatSize string  `json:"size_format"`
}

type IndexData struct {
	LoginDevices struct {
		Last struct {
			Device   string `json:"device"`
			DeviceId string `json:"device_id"`
			Ip       string `json:"ip"`
			City     string `json:"city"`
			Os       string `json:"os"`
			Network  string `json:"network"`
			Time     int64  `json:"utime"`
		} `json:"last"`
		List []LoginDeviceInfo `json:"list"`
	} `json:"login_devices_info"`
	Space struct {
		Total  SizeInfo `json:"all_total"`
		Used   SizeInfo `json:"all_use"`
		Remain SizeInfo `json:"all_remain"`
	} `json:"space_info"`
}
