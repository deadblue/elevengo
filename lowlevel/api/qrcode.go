package api

import (
	"fmt"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
)

const (
	qrcodeImageUrl = "https://qrcodeapi.115.com/api/1.0/web/1.0/qrcode?qrfrom=1&client=0&uid=%s"
)

type QrcodeTokenSpec struct {
	_JsonApiSpec[types.QrcodeTokenResult, protocol.QrcodeBaseResp]
}

func (s *QrcodeTokenSpec) Init() *QrcodeTokenSpec {
	s._JsonApiSpec.Init("https://qrcodeapi.115.com/api/1.0/web/1.0/token")
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

func (s *QrcodeLoginSpec) Init(uid string) *QrcodeLoginSpec {
	s._JsonApiSpec.Init("https://passportapi.115.com/app/1.0/web/1.0/login/qrcode")
	s.form.Set("account", uid).
		Set("app", "web")
	return s
}

func QrcodeImageUrl(userId string) string {
	return fmt.Sprintf(qrcodeImageUrl, userId)
}