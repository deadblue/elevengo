package webapi

type QrcodeTokenData struct {
	Uid    string `json:"uid"`
	Time   int64  `json:"time"`
	Sign   string `json:"sign"`
	Qrcode string `json:"qrcode"`
}

type QrcodeStatusData struct {
	Status  int    `json:"status,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Version string `json:"version,omitempty"`
}
