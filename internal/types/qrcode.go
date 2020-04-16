package types

import "encoding/json"

type QrcodeApiResult struct {
	State   int             `json:"state"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (r *QrcodeApiResult) IsFailed() bool {
	return r.State == 0
}

type QrcodeTokenData struct {
	Uid    string `json:"uid"`
	Time   int64  `json:"time"`
	Sign   string `json:"sign"`
	Qrcode string `json:"qrcode"`
}

type QrcodeStatusData struct {
	Status  int    `json:"status"`
	Msg     string `json:"msg"`
	Version string `json:"version"`
}

type QrcodeLoginData struct {
	Cookie struct {
		CID  string `json:"CID"`
		SEID string `json:"SEID"`
		UID  string `json:"UID"`
	} `json:"cookie"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Country  string `json:"country"`
	From     string `json:"from"`
}
