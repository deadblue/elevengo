package elevengo

import (
	"fmt"
	"time"
)

const (
	urlCaptchaApi = "https://captchaapi.115.com/"
)

func (c *Client) requestCaptchaImage(t string, index int) ([]byte, error) {
	qs := newRequestParameters().
		With("ct", "index").
		With("ac", "code")
	if t == "all" || t == "single" {
		qs.With("t", t)
		if t == "single" {
			qs.WithInt("id", index)
		}
	}
	qs.WithInt64("_t", time.Now().Unix())
	return c.requestRaw(urlCaptchaApi, qs, nil)
}

func (c *Client) CaptchaStart(client *string) (session *CaptchaSession, err error) {
	cb := fmt.Sprintf("Close911_%d", time.Now().UnixNano())
	// access captcha page
	qs := newRequestParameters().
		With("ac", "security_code").
		With("type", "web").
		With("cb", cb)
	if client != nil {
		qs.With("client", *client)
	}
	_, err = c.requestRaw(urlCaptchaApi, qs, nil)
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
