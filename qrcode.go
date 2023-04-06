package elevengo

import (
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"time"
)

// QrcodeSession holds the information during a QRCode login process.
type QrcodeSession struct {
	// The raw data of QRCode, caller should use third-party tools/libraries
	// to convert it into QRCode matrix or image.
	Content string
	// Hidden fields
	uid  string
	time int64
	sign string
}

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

func (a *Agent) qrcodeCallApi(url string, qs web.Params, form web.Payload, data interface{}) (err error) {
	resp := &webapi.LoginBasicResponse{}
	if err = a.wc.CallJsonApi(url, qs, form, resp); err != nil {
		return err
	}
	return resp.Decode(data)
}

// QrcodeStart starts a QRCode login session.
func (a *Agent) QrcodeStart(session *QrcodeSession) (err error) {
	data := &webapi.QrcodeTokenData{}
	if err = a.qrcodeCallApi(webapi.ApiQrcodeToken, nil, nil, data); err == nil {
		session.uid = data.Uid
		session.time = data.Time
		session.sign = data.Sign
		session.Content = data.Qrcode
	}
	return
}

// QrcodeStartForLinux starts a QRCode Img login for linux session.
// need use QrcodeLoginForLinux get cookie
func (a *Agent) QrcodeStartForLinux(session *QrcodeSession) (err error) {
	data := &webapi.QrcodeTokenSecretData{}
	now := time.Now().Unix()
	if err = a.wc.CallSecretJsonApi(webapi.ApiQrcodeTokenForLinux, nil, nil, data, now); err == nil {
		session.uid = data.Data.UID
		session.time = data.Data.Time
		session.sign = data.Data.Sign
		session.Content = webapi.ApiQrcodeImgForLinux + data.Data.UID
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
	qs := web.Params{}.
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
	form := web.Params{}.
		With("account", session.uid).
		With("app", "web").
		ToForm()
	data := &webapi.LoginUserData{}
	if err = a.qrcodeCallApi(webapi.ApiQrcodeLogin, nil, form, data); err == nil {
		a.uid = data.Id
	}
	return
}

// QrcodeLoginForLinux logins user through Linux QRCode.
func (a *Agent) QrcodeLoginForLinux(session *QrcodeSession) (err error) {
	form := web.Params{}.
		With("account", session.uid).
		With("app", "web").
		ToForm()
	data := &webapi.LoginUserData{}
	if err = a.qrcodeCallApi(webapi.ApiQrcodeLoginForLinux, nil, form, data); err == nil {
		a.uid = data.Id
	}
	return
}
