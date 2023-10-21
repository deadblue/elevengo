package api

import (
	"fmt"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
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

type QrcodeTokenSpec struct {
	_JsonApiSpec[types.QrcodeTokenResult, protocol.QrcodeBaseResp]
}

func (s *QrcodeTokenSpec) Init(appType string) *QrcodeTokenSpec {
	baseUrl := fmt.Sprintf(qrcodeTokenBaseUrl, appType)
	s._JsonApiSpec.Init(baseUrl)
	return s
}

type QrcodeStatusSpec struct {
	_JsonApiSpec[types.QrcodeStatusResult, protocol.QrcodeBaseResp]
}

func (s *QrcodeStatusSpec) Init(uid string, time int64, sign string) *QrcodeStatusSpec {
	s._JsonApiSpec.Init("https://qrcodeapi.115.com/get/status/")
	s.query.Set("uid", uid).
		SetInt64("time", time).
		Set("sign", sign).
		SetNow("_")
	return s
}

type QrcodeLoginSpec struct {
	_JsonApiSpec[types.QrcodeLoginResult, protocol.QrcodeBaseResp]
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
