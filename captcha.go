package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/util"
	"math/rand"
	"time"
)

func (c *Client) captchaValueImage() ([]byte, error) {
	qs := util.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithTimestamp("_t")
	return c.request(apiCaptcha, qs, nil)
}

func (c *Client) captchaKeyImage(index int) ([]byte, error) {
	qs := util.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithTimestamp("_t")
	if index < 0 || index > 9 {
		qs.WithString("t", "all")
	} else {
		qs.WithString("t", "single").WithInt("id", index)
	}
	return c.request(apiCaptcha, qs, nil)
}

func (c *Client) CaptchaStart() (session *CaptchaSession, err error) {
	cb := fmt.Sprintf("Close911_%d", time.Now().UnixNano())
	// request captcha page to start session
	qs := util.NewQueryString().
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("cb", cb)
	_, err = c.request(apiCaptcha, qs, nil)
	if err != nil {
		return
	}
	// request captcha images
	codeValue, err := c.captchaValueImage()
	if err != nil {
		return
	}
	codeKeys, err := c.captchaKeyImage(-1)
	if err != nil {
		return
	}
	// build session
	return &CaptchaSession{
		Callback:  cb,
		CodeValue: codeValue,
		CodeKeys:  codeKeys,
	}, nil
}

func (c *Client) CaptchaSubmit(code string, session *CaptchaSession) (err error) {
	// get captcha sign
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().UnixNano())
	qs := util.NewQueryString().
		WithString("ac", "code").
		WithString("t", "sign").
		WithString("callback", cb).
		WithTimestamp("_")
	signResult := &_CaptchaSignResult{}
	if err = c.requestJsonp(apiCaptcha, qs, signResult); err != nil {
		return
	}
	// post captcha code
	form := util.NewForm(false).
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("code", code).
		WithString("sign", signResult.Sign).
		WithString("cb", session.Callback)
	submitResult := &_BasicResult{}
	err = c.requestJson(apiCaptchaSubmit, nil, form, submitResult)
	if err == nil && !submitResult.State {
		err = ErrCaptchaFailed
	}
	return
}
