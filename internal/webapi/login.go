package webapi

import "encoding/json"

type LoginBasicResponse struct {
	State   int    `json:"state"`
	Code    int    `json:"code"`
	Message string `json:"message"`

	ErrorCode    int    `json:"errno"`
	ErrorMessage string `json:"error"`

	Data json.RawMessage `json:"data"`
}

func (r *LoginBasicResponse) Err() error {
	if r.State != 0 {
		return nil
	}
	return getError(r.Code)
}

func (r *LoginBasicResponse) Decode(data interface{}) error {
	return json.Unmarshal(r.Data, data)
}

type LoginKeyData struct {
	Key string `json:"key"`
}

type LoginResultData struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
	Mobile string `json:"mobile"`
}

type LoginUserData struct {
	Id     int    `json:"user_id"`
	Name   string `json:"user_name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	Cookie struct {
		CID  string `json:"CID"`
		UID  string `json:"UID"`
		SEID string `json:"SEID"`
	} `json:"cookie"`
}

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
