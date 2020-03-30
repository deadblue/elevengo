package elevengo

import (
	"encoding/json"
	"github.com/deadblue/elevengo/core"
	"strconv"
	"time"
)

type _QrcodeGenericResult struct {
	State   int             `json:"state"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type QrcodeToken struct {
	Uid     string `json:"uid"`
	Time    int64  `json:"time"`
	Sign    string `json:"sign"`
	Content string `json:"qrcode"`
}

type _QrcodeStatus struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
	Version string `json:"version"`
}

type _QrcodeLoginUser struct {
	Cookie   *Credentials `json:"cookie"`
	Country  string       `json:"country"`
	UserId   int          `json:"user_id"`
	UserName string       `json:"user_name"`
	Email    string       `json:"email"`
	Mobile   string       `json:"mobile"`
}

func (c *Client) callQrcodeApi(url string, qs *core.QueryString, form *core.Form, data interface{}) (err error) {
	result := new(_QrcodeGenericResult)
	err = c.requestJson(url, nil, nil, result)
	if err == nil && result.State != 1 {
		err = apiError(result.Code)
	}
	if err != nil || data == nil {
		return
	}
	return json.Unmarshal(result.Data, data)
}

func (c *Client) QrcodeStart() (session *QrcodeToken, err error) {
	token := new(QrcodeToken)
	if err = c.callQrcodeApi(apiQrcodeToken, nil, nil, token); err != nil {
		token = nil
	}
	return
}

func (c *Client) QrcodeCheck(token *QrcodeToken) (err error) {
	qs := core.NewQueryString().
		WithString("uid", token.Uid).
		WithString("time", strconv.FormatInt(token.Time, 10)).
		WithString("sign", token.Sign).
		WithInt64("_", time.Now().Unix())
	status := new(_QrcodeStatus)
	if err = c.callQrcodeApi(apiQrcodeStatus, qs, nil, status); err != nil {
		status = nil
		return
	}
	return
}

func (c *Client) QrcodeLogin(token *QrcodeToken) (err error) {
	form := core.NewForm(false).
		WithString("account", token.Uid).
		WithString("app", "web")
	result := new(_QrcodeLoginUser)
	err = c.callQrcodeApi(apiQrcodeLogin, nil, form, result)
	return
}
