package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/gostream/quietly"
	"io"
	"math/rand"
	"time"
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

// CaptchaStart starts a CAPTCHA session.
func (a *Agent) CaptchaStart(session *CaptchaSession) (err error) {
	// Fetch captcha page
	session.callback = fmt.Sprintf("Close911_%d", time.Now().UnixNano())
	qs := web.Params{}.With("cb", session.callback)
	body, err := a.wc.Get(webapi.ApiCaptchaPage, qs)
	if err != nil {
		return
	}
	util.ConsumeReader(body)

	// Fetch CAPTCHA code image
	qs = web.Params{}.WithNow("_t")
	if body, err = a.wc.Get(webapi.ApiCaptchaCodeImage, qs); err != nil {
		return
	}
	defer quietly.Close(body)
	if session.CodeImage, err = io.ReadAll(body); err != nil {
		return
	}

	// Fetch CAPTCHA keys image
	body, err = a.wc.Get(webapi.ApiCaptchaAllKeyImage, qs)
	if err != nil {
		return
	}
	defer quietly.Close(body)
	session.KeysImage, err = io.ReadAll(body)
	return
}

/*
CaptchaKeyImage gets one CAPTCHA key image data.

You can call this method multiple times, it will return the same character
in different font on every calling.

It is useful when you try to train your CAPTCHA solver.
*/
func (a *Agent) CaptchaKeyImage(session *CaptchaSession, index int) (data []byte, err error) {
	if index < 0 {
		index = 0
	} else if index > 9 {
		index = 9
	}
	qs := web.Params{}.
		WithInt("id", index).
		WithNow("_t")
	body, err := a.wc.Get(webapi.ApiCaptchaOneKeyImage, qs)
	if err != nil {
		return
	}
	defer quietly.Close(body)
	return io.ReadAll(body)
}

// CaptchaSubmit submits the CAPTCHA code to session.
func (a *Agent) CaptchaSubmit(session *CaptchaSession, code string) (err error) {
	// Get captcha sign
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().UnixNano())
	qs := web.Params{}.
		With("callback", cb).
		WithNow("_")
	signResp := &webapi.CaptchaSignResponse{}
	if err = a.wc.CallJsonpApi(webapi.ApiCaptchaSign, qs, signResp); err != nil {
		return
	}
	if err = signResp.Err(); err != nil {
		return
	}
	// Submit captcha code
	form := web.Params{}.
		With("ac", "security_code").
		With("type", "web").
		With("sign", signResp.Sign).
		With("code", code).
		With("cb", session.callback)
	submitResp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiCaptchaSubmit, nil, form, submitResp); err == nil {
		err = submitResp.Err()
	}
	return
}
