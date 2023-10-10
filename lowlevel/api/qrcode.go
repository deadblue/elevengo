package api

import (
	"encoding/json"
	"fmt"

	"github.com/deadblue/elevengo/lowlevel/errors"
)

const (
	qrcodeTokenBaseUrl = "https://qrcodeapi.115.com/api/1.0/%s/1.0/token"
	qrcodeLoginBaseUrl = "https://passportapi.115.com/app/1.0/%s/1.0/login/qrcode"
	qrcodeImageUrl     = "https://qrcodeapi.115.com/api/1.0/%s/1.0/qrcode?qrfrom=1&client=%d&uid=%s"
)

var (
	qrcodeAppIds = map[string]int{
		"web": 0,
		// Client ID for app is always 7
		"mac":     7,
		"linux":   7,
		"windows": 7,
	}
)

//lint:ignore U1000 This type is used in generic.
type _QrcodeBaseResp struct {
	State         int    `json:"state"`
	ErrorCode1    int    `json:"code"`
	ErrorCode2    int    `json:"errno"`
	ErrorMessage1 string `json:"message"`
	ErrorMessage2 string `json:"error"`

	Data json.RawMessage `json:"data"`
}

func (r *_QrcodeBaseResp) Err() error {
	if r.State != 0 {
		return nil
	}
	return errors.Get(r.ErrorCode1)
}

func (r *_QrcodeBaseResp) Extract(v any) error {
	return json.Unmarshal(r.Data, v)
}

type QrcodeTokenResult struct {
	Uid  string `json:"uid"`
	Time int64  `json:"time"`
	Sign string `json:"sign"`
}

type QrcodeTokenSpec struct {
	_JsonApiSpec[QrcodeTokenResult, _QrcodeBaseResp]
}

func (s *QrcodeTokenSpec) Init(appType string) *QrcodeTokenSpec {
	baseUrl := fmt.Sprintf(qrcodeTokenBaseUrl, appType)
	s._JsonApiSpec.Init(baseUrl)
	return s
}

type QrcodeStatusResult struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"msg,omitempty"`
	Version string `json:"version,omitempty"`
}

type QrcodeStatusSpec struct {
	_JsonApiSpec[QrcodeStatusResult, _QrcodeBaseResp]
}

func (s *QrcodeStatusSpec) Init(uid string, time int64, sign string) *QrcodeStatusSpec {
	s._JsonApiSpec.Init("https://qrcodeapi.115.com/get/status/")
	s.query.Set("uid", uid).
		SetInt64("time", time).
		Set("sign", sign).
		SetNow("_")
	return s
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

type QrcodeLoginSpec struct {
	_JsonApiSpec[QrcodeLoginResult, _QrcodeBaseResp]
}

func (s *QrcodeLoginSpec) Init(appType string, uid string) *QrcodeLoginSpec {
	baseUrl := fmt.Sprintf(qrcodeLoginBaseUrl, appType)
	s._JsonApiSpec.Init(baseUrl)
	s.form.Set("account", uid).
		Set("app", appType)
	return s
}

func QrcodeImageUrl(appType, userId string) string {
	appId := qrcodeAppIds[appType]
	return fmt.Sprintf(qrcodeImageUrl, appType, appId, userId)
}
