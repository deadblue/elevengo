package elevengo

import (
	"errors"
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"math/rand"
	"time"
)

const (
	apiCaptcha       = "https://captchaapi.115.com/"
	apiCaptchaSubmit = "https://webapi.115.com/user/captcha"
)

type CaptchaSession struct {
	callback string
	code     []byte
}

func (cs *CaptchaSession) CodeImage() []byte {
	return cs.code
}

func (c *Client) CaptchaStart() (session *CaptchaSession, err error) {
	// Fetch captcha page
	callback := fmt.Sprintf("Close911_%d", time.Now().UnixNano())
	qs := core.NewQueryString().
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("cb", callback)
	if _, err = c.hc.Get(apiCaptcha, qs); err != nil {
		return
	}
	// Fetch captcha value image
	qs = core.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithInt64("_t", time.Now().Unix())
	code, err := c.hc.Get(apiCaptcha, qs)
	if err != nil {
		return
	}
	// Build session
	session = &CaptchaSession{
		callback: callback,
		code:     code,
	}
	return
}

func (c *Client) CaptchaAllKeysImage(session *CaptchaSession) ([]byte, error) {
	qs := core.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithString("t", "all").
		WithInt64("_t", time.Now().Unix())
	return c.hc.Get(apiCaptcha, qs)
}

func (c *Client) CaptchaKeyImage(session *CaptchaSession, index int) ([]byte, error) {
	if index < 0 {
		index = 0
	} else if index > 9 {
		index = 9
	}
	qs := core.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithString("t", "single").
		WithInt("id", index).
		WithInt64("_t", time.Now().Unix())
	return c.hc.Get(apiCaptcha, qs)
}

func (c *Client) CaptchaSubmit(session *CaptchaSession, code string) (err error) {
	// Get captcha sign
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().UnixNano())
	qs := core.NewQueryString().
		WithString("ac", "code").
		WithString("t", "sign").
		WithString("callback", cb).
		WithInt64("_", time.Now().Unix())
	signResult := &internal.CaptchaSignResult{}
	if err = c.hc.JsonpApi(apiCaptcha, qs, signResult); err != nil {
		return
	}
	// Submit captcha code
	form := core.NewForm().
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("sign", signResult.Sign).
		WithString("code", code).
		WithString("cb", session.callback)
	submitResult := &internal.CaptchaSubmitResult{}
	err = c.hc.JsonApi(apiCaptchaSubmit, nil, form, submitResult)
	if err == nil && !submitResult.State {
		// TODO: handle submit result
		err = errors.New("submit failed")
	}
	return
}
