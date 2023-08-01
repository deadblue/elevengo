package api

import (
	"fmt"

	"github.com/deadblue/elevengo/internal/api/base"
)

type _QrcodeBaseResp[D any] struct {
	State         int    `json:"state"`
	ErrorCode1    int    `json:"code"`
	ErrorCode2    int    `json:"errno"`
	ErrorMessage1 string `json:"message"`
	ErrorMessage2 string `json:"error"`

	Data D `json:"data"`
}

func (r _QrcodeBaseResp[D]) Err() error {
	if r.State != 0 {
		return nil
	}
	return base.GetError(r.ErrorCode1)
}

type _QrcodeTokenData struct {
	Uid  string `json:"uid"`
	Time int64  `json:"time"`
	Sign string `json:"sign"`
}

type _QrcodeStatusData struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"msg,omitempty"`
	Version string `json:"version,omitempty"`
}

type _QrcodeLoginData struct {
	Cookie struct {
		CID  string `json:"CID"`
		SEID string `json:"SEID"`
		UID  string `json:"UID"`
	} `json:"cookie"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
}

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

type QrcodeTokenSpec struct {
	base.JsonApiSpec[_QrcodeBaseResp[_QrcodeTokenData]]
}

func (s *QrcodeTokenSpec) Init(appType string) *QrcodeTokenSpec {
	baseUrl := fmt.Sprintf(qrcodeTokenBaseUrl, appType)
	s.JsonApiSpec.Init(baseUrl)
	return s
}

type QrcodeStatusSpec struct {
	base.JsonApiSpec[_QrcodeBaseResp[_QrcodeStatusData]]
}

func (s *QrcodeStatusSpec) Init(uid string, time int64, sign string) *QrcodeStatusSpec {
	s.JsonApiSpec.Init("https://qrcodeapi.115.com/get/status/")
	s.QuerySet("uid", uid)
	s.QuerySetInt64("time", time)
	s.QuerySet("sign", sign)
	s.QuerySetNow("_")
	return s
}

type QrcodeLoginSpec struct {
	base.JsonApiSpec[_QrcodeBaseResp[_QrcodeLoginData]]
}

func (s *QrcodeLoginSpec) Init(appType string, uid string) *QrcodeLoginSpec {
	baseUrl := fmt.Sprintf(qrcodeLoginBaseUrl, appType)
	s.JsonApiSpec.Init(baseUrl)
	s.FormSet("account", uid)
	s.FormSet("app", appType)
	return s
}

func QrcodeImageUrl(appType, userId string) string {
	appId := qrcodeAppIds[appType]
	return fmt.Sprintf(qrcodeImageUrl, appType, appId, userId)
}
