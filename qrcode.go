package elevengo

import (
	"encoding/json"
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"strconv"
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
	uid  string
	time int64
	sign string
	// The raw data of QRcode, you should use thridparty
	// tools/libraries to convert it into QRcode image.
	Content []byte
}

// QrcodeStatus is returned by `Client.QrcodeStatus()`.
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

// Return true if user allowed this login process, you can call "Client.QrcodeLogin()" after then.
func (qs QrcodeStatus) IsAllowed() bool {
	return qs == 2
}

// Return true if user canceled this login process.
func (qs QrcodeStatus) IsCanceled() bool {
	return qs == -2
}

// Return true if the QRcode is expired.
func (qs QrcodeStatus) IsExpired() bool {
	return qs == codeQrcodeExpired
}

type qrcodeError struct {
	code int
}

func (qe *qrcodeError) Error() string {
	return fmt.Sprintf("upstream qrcode API error: %d", qe.code)
}

func (c *Client) callQrcodeApi(url string, qs core.QueryString, form core.Form, data interface{}) error {
	result := &internal.QrcodeApiResult{}
	if err := c.hc.JsonApi(url, qs, form, result); err != nil {
		return err
	}
	if result.IsFailed() {
		return &qrcodeError{
			code: result.Code,
		}
	}
	return json.Unmarshal(result.Data, data)
}

// Start a QRcode login process.
func (c *Client) QrcodeStart() (session *QrcodeSession, err error) {
	data := &internal.QrcodeTokenData{}
	if err = c.callQrcodeApi(apiQrcodeToken, nil, nil, data); err == nil {
		session = &QrcodeSession{
			uid:     data.Uid,
			time:    data.Time,
			sign:    data.Sign,
			Content: []byte(data.Qrcode),
		}
	}
	return
}

// Get QRcode login process status.
// The remote API uses a long-pull request for 30 seconds, so this API
// will also block at most 30 seconds, be careful to use it in main goroutine.
func (c *Client) QrcodeStatus(session *QrcodeSession) (status QrcodeStatus, err error) {
	qs := core.NewQueryString().
		WithString("uid", session.uid).
		WithInt64("time", session.time).
		WithString("sign", session.sign).
		WithInt64("_", time.Now().Unix())
	data := &internal.QrcodeStatusData{}
	if err = c.callQrcodeApi(apiQrcodeStatus, qs, nil, data); err == nil {
		status = QrcodeStatus(data.Status)
	} else {
		if qerr, ok := err.(*qrcodeError); ok && qerr.code == codeQrcodeExpired {
			status, err = QrcodeStatus(codeQrcodeExpired), nil
		}
	}
	return
}

// Login through QRcode.
// You SHOULD call this method ONLY when `QrcodeStatus.IsAllowed()` is true.
func (c *Client) QrcodeLogin(session *QrcodeSession) error {
	form := core.NewForm().
		WithString("account", session.uid).
		WithString("app", "web")
	data := &internal.QrcodeLoginData{}
	if err := c.callQrcodeApi(apiQrcodeLogin, nil, form, data); err != nil {
		return err
	} else {
		c.ui = &internal.UserInfo{
			UserId: strconv.Itoa(data.UserId),
		}
		return nil
	}
}
