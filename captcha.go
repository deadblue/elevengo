package elevengo

import (
	"fmt"
	"time"
)

const (
	urlCaptchaApi = "https://captchaapi.115.com/"
)

func (c *Client) requestCaptchaImage(t string, index int) ([]byte, error) {
	qs := newQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithTimestamp("_t")
	if t == "all" || t == "single" {
		qs.WithString("t", t)
		if t == "single" {
			qs.WithInt("id", index)
		}
	}
	return c.request(urlCaptchaApi, qs, nil)
}

func (c *Client) CaptchaStart(client *string) (session *CaptchaSession, err error) {
	cb := fmt.Sprintf("Close911_%d", time.Now().UnixNano())
	// access captcha page
	qs := newQueryString().
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("cb", cb)
	if client != nil {
		qs.WithString("client", *client)
	}
	_, err = c.request(urlCaptchaApi, qs, nil)
	if err != nil {
		return
	}
	// request code value image
	codeValue, err := c.requestCaptchaImage("", 0)
	if err != nil {
		return
	}
	codeKeys, err := c.requestCaptchaImage("all", 0)
	if err != nil {
		return
	}
	session = &CaptchaSession{
		Callback:  cb,
		CodeValue: codeValue,
		CodeKeys:  codeKeys,
	}
	return
}

func (c *Client) CaptchaSubmit(code string, session *CaptchaSession) (err error) {

	return nil
}
