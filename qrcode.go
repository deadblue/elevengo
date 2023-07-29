package elevengo

import (
	"errors"
	"io"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/elevengo/option"
)

// QrcodeSession holds the information during a QRCode login process.
type QrcodeSession struct {
	// QRCode image content.
	Image []byte
	// Hidden fields.
	uid  string
	time int64
	sign string
	appType string
}

var ErrQrcodeCancelled = errors.New("QRcode cancelled")

func (a *Agent) qrcodeCallApi(url string, qs protocol.Params, form protocol.Payload, data interface{}) (err error) {
	resp := &webapi.LoginBasicResponse{}
	if err = a.pc.CallJsonApi(url, qs, form, resp); err != nil {
		return err
	}
	return resp.Decode(data)
}

// QrcodeStart starts a QRcode sign-in session.
// The session is for web by default, you can change sign-in app by passing a 
// "option.QrcodeLoginOption".
// 
// Example:
// 
//   agent := elevengo.Default()
//   session := elevengo.QrcodeSession()
//   agent.QrcodeStart(session, option.QrcodeLoginLinux)
// 
func (a *Agent) QrcodeStart(session *QrcodeSession, options ...option.QrcodeOption) (err error) {
	// Apply options
	for _, opt := range options  {
		switch opt := opt.(type) {
		case option.QrcodeLoginOption:
			session.appType = string(opt)
		}
	}
	if session.appType == "" {
		session.appType = string(option.QrcodeLoginWeb)
	}
	// Start session
	data := &webapi.QrcodeTokenData{}
	if err = a.qrcodeCallApi(webapi.QrcodeTokenApi(session.appType), nil, nil, data); err != nil {
		return
	}
	session.uid = data.Uid
	session.time = data.Time
	session.sign = data.Sign
	// Fetch QRcode image data
	var reader io.ReadCloser
	if reader, err = a.Get(webapi.QrcodeImageUrl(session.appType, data.Uid)); err != nil {
		return
	}
	defer util.QuietlyClose(reader)
	session.Image, err = io.ReadAll(reader)
	return
}

func (a *Agent) qrcodeSignIn(session *QrcodeSession) (err error) {
	form := protocol.Params{}.
		With("account", session.uid).
		With("app", session.appType).
		ToForm()
	data := &webapi.LoginUserData{}
	return a.qrcodeCallApi(
		webapi.QrcodeLoginApi(session.appType), nil, form, data,
	)
}

// QrcodePoll polls the session state, and automatically sin
func (a *Agent) QrcodePoll(session *QrcodeSession) (done bool, err error) {
	qs := protocol.Params{}.
		With("uid", session.uid).
		WithInt64("time", session.time).
		With("sign", session.sign).
		WithNow("_")
	data := &webapi.QrcodeStatusData{}
	if err = a.qrcodeCallApi(webapi.ApiQrcodeStatus, qs, nil, data); err != nil {
		return
	}
	switch data.Status {
	case -2:
		err = ErrQrcodeCancelled
	case 2:
		err = a.qrcodeSignIn(session)
		done = err == nil
	}
	return 
}
