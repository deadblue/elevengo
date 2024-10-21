package elevengo

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/option"
)

const (
	formatQrcodeSession = "%s|%s|%d|%s"
)

// QrcodeSession holds the information during a QRCode login process.
type QrcodeSession struct {
	// QRCode image content.
	Image []byte
	// Hidden fields.
	app  string
	uid  string
	time int64
	sign string
}

// Marshal marshals QrcodeSession into a string, which can be tranfered out of
// process.
func (s *QrcodeSession) Marshal() string {
	if s.time == 0 {
		return ""
	}
	return fmt.Sprintf(formatQrcodeSession, s.app, s.uid, s.time, s.sign)
}

// Unmarshal fills QrcodeSession from a string which is generated by |Marshal|.
func (s *QrcodeSession) Unmarshal(text string) (err error) {
	_, err = fmt.Sscanf(
		text, formatQrcodeSession,
		&s.app, &s.uid, &s.time, &s.sign,
	)
	return
}

var ErrQrcodeCancelled = errors.New("QRcode cancelled")

// QrcodeStart starts a QRcode sign-in session.
// The session is for web by default, you can change sign-in app by passing a
// "option.QrcodeLoginOption".
//
// Example:
//
//	agent := elevengo.Default()
//	session := elevengo.QrcodeSession()
//	agent.QrcodeStart(session, option.Qrcode().LoginTv())
func (a *Agent) QrcodeStart(session *QrcodeSession, options ...*option.QrcodeOptions) (err error) {
	// Apply options
	app := "web"
	if opts := util.NotNull(options...); opts != nil {
		app = opts.App
	}
	// Get token
	spec := (&api.QrcodeTokenSpec{}).Init(app)
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}
	session.app = app
	session.uid = spec.Result.Uid
	session.time = spec.Result.Time
	session.sign = spec.Result.Sign
	// Fetch QRcode image data
	var reader io.ReadCloser
	if reader, err = a.Fetch(api.QrcodeImageUrl(session.app, session.uid)); err != nil {
		return
	}
	defer util.QuietlyClose(reader)
	session.Image, err = io.ReadAll(reader)
	return
}

func (a *Agent) qrcodeSignIn(session *QrcodeSession) (err error) {
	spec := (&api.QrcodeLoginSpec{}).Init(session.app, session.uid)
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}
	return a.afterSignIn()
}

// QrcodePoll polls the session state, and automatically sin
func (a *Agent) QrcodePoll(session *QrcodeSession) (done bool, err error) {
	spec := (&api.QrcodeStatusSpec{}).Init(
		session.uid, session.time, session.sign,
	)
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}
	switch spec.Result.Status {
	case -2:
		err = ErrQrcodeCancelled
	case 2:
		err = a.qrcodeSignIn(session)
		done = err == nil
	}
	return
}
