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
		UID  string `json:"UID"`
		CID  string `json:"CID"`
		KID  string `json:"KID"`
		SEID string `json:"SEID"`
	} `json:"cookie"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}
