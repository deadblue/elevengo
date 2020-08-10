package elevengo

import (
	"encoding/json"
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"time"
)

const (
	apiQrcodeToken  = "https://qrcodeapi.115.com/api/1.0/web/1.0/token"
	apiQrcodeStatus = "https://qrcodeapi.115.com/get/status/"
	apiQrcodeLogin  = "https://passportapi.115.com/app/1.0/web/1.0/login/qrcode"
)

// QrcodeSession holds the information during a QRcode login process.
type QrcodeSession struct {
	// The raw data of QR code, caller should use thridparty tools/libraries
	// to convert it into QR code matrix or image.
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

// Return true if user still does not scan the QRcode.
func (qs QrcodeStatus) IsWaiting() bool {
	return qs == 0
}

// Return true if user has scanned the QRcode, but still not allow or cancel this login.
func (qs QrcodeStatus) IsScanned() bool {
	return qs == 1
}

// Return true if user allowed this login process, you can call "Agent.QrcodeLogin()" after then.
func (qs QrcodeStatus) IsAllowed() bool {
	return qs == 2
}

// Return true if user canceled this login process.
func (qs QrcodeStatus) IsCanceled() bool {
	return qs == -2
}

func (a *Agent) callQrcodeApi(url string, qs core.QueryString, form core.Form, data interface{}) error {
	result := &types.QrcodeApiResult{}
	if err := a.hc.JsonApi(url, qs, form, result); err != nil {
		return err
	}
	if result.IsFailed() {
		return types.MakeQrcodeError(result.Code, result.Message)
	}
	return json.Unmarshal(result.Data, data)
}

// Start a QRcode login session.
func (a *Agent) QrcodeStart() (session *QrcodeSession, err error) {
	data := &types.QrcodeTokenData{}
	if err = a.callQrcodeApi(apiQrcodeToken, nil, nil, data); err == nil {
		session = &QrcodeSession{
			uid:     data.Uid,
			time:    data.Time,
			sign:    data.Sign,
			Content: data.Qrcode,
		}
	}
	return
}

/*
Get QRcode status.

The upstream API uses a long-pull request for 30 seconds, so this API will
also block at most 30 seconds, be careful to use it in main goroutine.

The QRcode has 4 status:

	- Waiting
	- Scanned
	- Allowed
	- Canceled

The QRcode will expire in 5 mimutes, when it expired, an error will be return, caller
can use IsQrcodeExipre() to check that.
*/
func (a *Agent) QrcodeStatus(session *QrcodeSession) (status QrcodeStatus, err error) {
	qs := core.NewQueryString().
		WithString("uid", session.uid).
		WithInt64("time", session.time).
		WithString("sign", session.sign).
		WithInt64("_", time.Now().Unix())
	data := &types.QrcodeStatusData{}
	if err = a.callQrcodeApi(apiQrcodeStatus, qs, nil, data); err == nil {
		status = QrcodeStatus(data.Status)
	}
	return
}

// Login through QRcode.
// You SHOULD call this method ONLY when `QrcodeStatus.IsAllowed()` is true.
func (a *Agent) QrcodeLogin(session *QrcodeSession) error {
	form := core.NewForm().
		WithString("account", session.uid).
		WithString("app", "web")
	data := &types.QrcodeLoginData{}
	if err := a.callQrcodeApi(apiQrcodeLogin, nil, form, data); err != nil {
		return err
	} else {
		a.ui = &UserInfo{
			Id:   data.UserId,
			Name: data.UserName,
		}
		return nil
	}
}
