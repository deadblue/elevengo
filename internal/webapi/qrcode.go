package webapi

import "fmt"

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

const (
	_ApiFormatQrcodeToken = "https://qrcodeapi.115.com/api/1.0/%s/1.0/token"
	_ApiFormatQrcodeLogin = "https://passportapi.115.com/app/1.0/%s/1.0/login/qrcode"
	_UrlFormatQrcodeImage = "https://qrcodeapi.115.com/api/1.0/%s/1.0/qrcode?qrfrom=1&client=%d&uid=%s"
)

var (
	_PlatformIdMap = map[string]int{
		"windows": 5,
		"linux": 7,
	}
)

func QrcodeTokenApi(platform string) string {
	return fmt.Sprintf(_ApiFormatQrcodeToken, platform)
}

func QrcodeLoginApi(platform string) string {
	return fmt.Sprintf(_ApiFormatQrcodeLogin, platform)
}

func QrcodeImageUrl(platform, userId string) string {
	platformId := _PlatformIdMap[platform]
	return fmt.Sprintf(_UrlFormatQrcodeImage, platform, platformId, userId)
}