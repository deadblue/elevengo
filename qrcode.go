package elevengo

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
)

// QrcodeSession holds the information during a QRCode login process.
type QrcodeSession struct {
	// The raw data of QRCode, caller should use third-party tools/libraries
	// to convert it into QRCode matrix or image.
	// 
	// Deprecated: Please use `ImageUrl` instead.
	Content string
	// URL of QRCode image.
	ImageUrl string
	// Hidden fields
	uid  string
	time int64
	sign string
	platform string
}

type QrcodePlatform string

const (
	QrcodePlatformLinux   QrcodePlatform = "linux"
	QrcodePlatformMac     QrcodePlatform = "mac"
	QrcodePlatformWindows QrcodePlatform = "windows"
	QrcodePlatformWeb     QrcodePlatform = "web"
)

// QrcodeStatus is returned by `Agent.QrcodeStatus()`.
// You can call `QrcodeStatus.IsXXX()` method to check the status,
// or directly check its value.
type QrcodeStatus int

func (s QrcodeStatus) IsWaiting() bool {
	return s == 0
}
func (s QrcodeStatus) IsScanned() bool {
	return s == 1
}
func (s QrcodeStatus) IsAllowed() bool {
	return s == 2
}
func (s QrcodeStatus) IsCanceled() bool {
	return s == -2
}

func (a *Agent) qrcodeCallApi(url string, qs protocol.Params, form protocol.Payload, data interface{}) (err error) {
	resp := &webapi.LoginBasicResponse{}
	if err = a.pc.CallJsonApi(url, qs, form, resp); err != nil {
		return err
	}
	return resp.Decode(data)
}

// QrcodeStart starts a QRCode login session for web.
func (a *Agent) QrcodeStart(session *QrcodeSession) (err error) {
	return a.QrcodeStartForPlatform(session, QrcodePlatformWeb)
}

// QrcodeStartForPlatform starts a QRCode login session for specific platform.
func (a *Agent) QrcodeStartForPlatform(session *QrcodeSession, platform QrcodePlatform) (err error) {
	data := &webapi.QrcodeTokenData{}
	if err = a.qrcodeCallApi(webapi.QrcodeTokenApi(string(platform)), nil, nil, data); err == nil {
		session.platform = string(platform)
		session.uid = data.Uid
		session.time = data.Time
		session.sign = data.Sign
		session.ImageUrl = webapi.QrcodeImageUrl(session.platform, data.Uid)
		if platform == QrcodePlatformWeb {
			session.Content = data.Qrcode
		}
	}
	return
}

/*
QrcodeStatus returns the status of QRCode login session.

The upstream API uses a long-pull request for 30 seconds, so this API will
also block at most 30 seconds, be careful to use it in main goroutine.

There will be 4 kinds of status:

  - Waiting
  - Scanned
  - Allowed
  - Canceled

The QRCode will expire in 5 minutes, when it expired, an error will be return, caller
can use IsQrcodeExpire() to check that.
*/
func (a *Agent) QrcodeStatus(session *QrcodeSession) (status QrcodeStatus, err error) {
	qs := protocol.Params{}.
		With("uid", session.uid).
		WithInt64("time", session.time).
		With("sign", session.sign).
		WithNow("_")
	data := &webapi.QrcodeStatusData{}
	if err = a.qrcodeCallApi(webapi.ApiQrcodeStatus, qs, nil, data); err == nil {
		status = QrcodeStatus(data.Status)
	}
	return
}

// QrcodeLogin logins user through QRCode.
// You SHOULD call this method ONLY when `QrcodeStatus.IsAllowed()` is true.
func (a *Agent) QrcodeLogin(session *QrcodeSession) (err error) {
	form := protocol.Params{}.
		With("account", session.uid).
		With("app", session.platform).
		ToForm()
	data := &webapi.LoginUserData{}
	return a.qrcodeCallApi(
		webapi.QrcodeLoginApi(session.platform), nil, form, data,
	)
}
