package elevengo

import (
	"encoding/json"
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"io"
	"strings"
	"time"
)

const (
	apiQrcodeToken  = "https://qrcodeapi.115.com/api/1.0/web/1.0/token"
	apiQrcodeStatus = "https://qrcodeapi.115.com/get/status/"
	apiQrcodeLogin  = "https://passportapi.115.com/app/1.0/web/1.0/login/qrcode"

	codeQrcodeExpired = 40199002
)

// QrcodeSession holds the information during a QRcode login process.
type QrcodeSession struct {
	uid     string
	time    int64
	sign    string
	content string
}

// Get the raw data of QRcode.
// You should use a thridparty tools/libraries to convert it into QRcode image.
func (qs *QrcodeSession) Content() io.Reader {
	return strings.NewReader(qs.content)
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

type QrcodeError struct {
	code int
}

func (qe *QrcodeError) Error() string {
	return fmt.Sprintf("upstream qrcode API error: %d", qe.code)
}

func (qe *QrcodeError) IsExpired() bool {
	return qe.code == codeQrcodeExpired
}

func (a *Agent) callQrcodeApi(url string, qs core.QueryString, form core.Form, data interface{}) error {
	result := &internal.QrcodeApiResult{}
	if err := a.hc.JsonApi(url, qs, form, result); err != nil {
		return err
	}
	if result.IsFailed() {
		return &QrcodeError{
			code: result.Code,
		}
	}
	return json.Unmarshal(result.Data, data)
}

// Start a QRcode login process.
func (a *Agent) QrcodeStart() (session *QrcodeSession, err error) {
	data := &internal.QrcodeTokenData{}
	if err = a.callQrcodeApi(apiQrcodeToken, nil, nil, data); err == nil {
		session = &QrcodeSession{
			uid:     data.Uid,
			time:    data.Time,
			sign:    data.Sign,
			content: data.Qrcode,
		}
	}
	return
}

// Get QRcode login process status.
// The remote API uses a long-pull request for 30 seconds, so this API
// will also block at most 30 seconds, be careful to use it in main goroutine.
func (a *Agent) QrcodeStatus(session *QrcodeSession) (status QrcodeStatus, err error) {
	qs := core.NewQueryString().
		WithString("uid", session.uid).
		WithInt64("time", session.time).
		WithString("sign", session.sign).
		WithInt64("_", time.Now().Unix())
	data := &internal.QrcodeStatusData{}
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
	data := &internal.QrcodeLoginData{}
	if err := a.callQrcodeApi(apiQrcodeLogin, nil, form, data); err != nil {
		return err
	} else {
		a.ui = &internal.UserInfo{
			UserId: data.UserId,
		}
		return nil
	}
}
