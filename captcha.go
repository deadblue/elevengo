package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"math/rand"
	"time"
)

const (
	apiCaptcha       = "https://captchaapi.115.com/"
	apiCaptchaSubmit = "https://webapi.115.com/user/captcha"
)

/*
CaptchaSession holds CAPTCHA images and session information during a CAPTCHA process.

There are 4 Chinese characters on "CodeImage" and 10 Chinese characters on
"KeysImage" (5 columns by 2 rows). User should find the 4 characters on "CodeImage"
from "KeysImage", the indexes of these characters is the CAPTCHA code.

The index bases on 0, starts from left top corner, increases from left to right,
then from top to bottom.

For example:

Assume the code image is
	+---------------+
	| J | F | C | H |
	+---------------+
and the keys image is
	+-------------------+
	| A | B | C | D | E |
	+-------------------+
	| F | G | H | I | J |
	+-------------------+
Then the CAPTCHA code is "9527", you can call Agent.CaptchaSubmit() to submit it.
*/
type CaptchaSession struct {
	// CAPTCHA image data.
	CodeImage []byte

	// CAPTCHA keys image data.
	KeysImage []byte

	// Hidden fields.
	callback string
}

// Start a CAPTCHA session.
func (a *Agent) CaptchaStart() (session *CaptchaSession, err error) {
	// Fetch captcha page
	callback := fmt.Sprintf("Close911_%d", time.Now().UnixNano())
	qs := core.NewQueryString().
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("cb", callback)
	if _, err = a.hc.Get(apiCaptcha, qs); err != nil {
		return
	}
	// Fetch CAPTCHA code image
	qs = core.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithInt64("_t", time.Now().Unix())
	codeImg, err := a.hc.Get(apiCaptcha, qs)
	if err != nil {
		return
	}
	// Fetch CAPTCHA keys image
	qs = core.NewQueryString().
		WithString("ct", "index").
		WithString("ac", "code").
		WithString("t", "all").
		WithInt64("_t", time.Now().Unix())
	keysImg, err := a.hc.Get(apiCaptcha, qs)
	if err != nil {
		return
	}
	// Build session
	session = &CaptchaSession{
		callback:  callback,
		CodeImage: codeImg,
		KeysImage: keysImg,
	}
	return
}

/*
Get one CAPTCHA key image data.

You can call this method multiple times, it will return the same character
in different font on every calling.

It is useful when you try to train your CAPTCHA solver.
*/
func (a *Agent) CaptchaKeyImage(session *CaptchaSession, index int) ([]byte, error) {
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
	return a.hc.Get(apiCaptcha, qs)
}

// Submit the CAPTCHA code.
func (a *Agent) CaptchaSubmit(session *CaptchaSession, code string) (err error) {
	// Get captcha sign
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().UnixNano())
	qs := core.NewQueryString().
		WithString("ac", "code").
		WithString("t", "sign").
		WithString("callback", cb).
		WithInt64("_", time.Now().Unix())
	signResult := &types.CaptchaSignResult{}
	if err = a.hc.JsonpApi(apiCaptcha, qs, signResult); err != nil {
		return
	}
	// Submit captcha code
	form := core.NewForm().
		WithString("ac", "security_code").
		WithString("type", "web").
		WithString("sign", signResult.Sign).
		WithString("code", code).
		WithString("cb", session.callback)
	submitResult := &types.CaptchaSubmitResult{}
	err = a.hc.JsonApi(apiCaptchaSubmit, nil, form, submitResult)
	if err == nil && submitResult.IsFailed() {
		err = errCaptchaFailed
	}
	return
}
