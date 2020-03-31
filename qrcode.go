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

	apiLoginQrcode = "https://passportapi.115.com/app/1.0/web/1.0/login/qrcode"
)

// QrcodeSession holds the information during a QRcode login process.
type QrcodeSession struct {
	uid     string
	time    int64
	sign    string
	content []byte
}

// Get the raw data of qrcode.
// You should use a thridparty tools/libraries to convert it into QRcode image.
func (qs *QrcodeSession) Content() []byte {
	return qs.content
}

// QrcodeStatus is returned by `Client.QrcodeStatus()`.
// You can call `QrcodeStatus.IsXXX()` method to check the status,
// or directly check its value.
type QrcodeStatus int

func (qs QrcodeStatus) IsWaiting() bool {
	return qs == 0
}
func (qs QrcodeStatus) IsScanned() bool {
	return qs == 1
}
func (qs QrcodeStatus) IsAllowed() bool {
	return qs == 2
}
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
	return qe.code == 40199002
}

func (c *Client) callQrcodeApi(url string, qs core.QueryString, data interface{}) error {
	result := &internal.QrcodeApiResult{}
	if err := c.hc.JsonApi(url, qs, nil, result); err != nil {
		return err
	}
	if result.State != 1 {
		return &QrcodeError{
			code: result.Code,
		}
	}
	return json.Unmarshal(result.Data, data)
}

// Start a QRcode login process
func (c *Client) QrcodeStart() (session *QrcodeSession, err error) {
	data := &internal.QrcodeTokenData{}
	if err = c.callQrcodeApi(apiQrcodeToken, nil, data); err == nil {
		session = &QrcodeSession{
			uid:     data.Uid,
			time:    data.Time,
			sign:    data.Sign,
			content: []byte(data.Qrcode),
		}
	}
	return
}

// Get QRcode login session status.
//
// The upstream API uses a long-pull request for 30 seconds, so this API
// will also block at most 30 seconds, be careful to use it in main goroutine.
func (c *Client) QrcodeStatus(session *QrcodeSession) (status QrcodeStatus, err error) {
	qs := core.NewQueryString().
		WithString("uid", session.uid).
		WithInt64("time", session.time).
		WithString("sign", session.sign).
		WithInt64("_", time.Now().Unix())
	data := &internal.QrcodeStatusData{}
	if err = c.callQrcodeApi(apiQrcodeStatus, qs, data); err == nil {
		status = QrcodeStatus(data.Status)
	}
	return
}

// Login by QRcode.
//
// You should call this method ONLY when `QrcodeStatus.IsAllowed()` is true.
func (c *Client) LoginByQrcode(session *QrcodeSession) error {
	form := core.NewForm().
		WithString("account", session.uid).
		WithString("app", "web")
	result := &internal.QrcodeLoginResult{}
	if err := c.hc.JsonApi(apiLoginQrcode, nil, form, result); err != nil {
		return err
	} else {
		c.ui = &internal.UserInfo{
			UserId: strconv.Itoa(result.Data.UserId),
		}
		return nil
	}
}
