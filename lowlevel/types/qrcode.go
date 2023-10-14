package types

type QrcodeTokenResult struct {
	Uid  string `json:"uid"`
	Time int64  `json:"time"`
	Sign string `json:"sign"`
}

type QrcodeStatusResult struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"msg,omitempty"`
	Version string `json:"version,omitempty"`
}

type QrcodeLoginResult struct {
	Cookie struct {
		CID  string `json:"CID"`
		SEID string `json:"SEID"`
		UID  string `json:"UID"`
	} `json:"cookie"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}
