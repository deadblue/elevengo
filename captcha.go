package elevengo

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
)

// IsCaptchaRequired indicates whether err requires user to solve a captcha.
func IsCaptchaRequired(err error) bool {
	return errors.Is(err, webapi.ErrCaptchaRequired)
}

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
	session.callback = fmt.Sprintf("Close911_%d", time.Now().UnixNano())
	// Access captcha page
	qs := protocol.Params{}.With("cb", session.callback)
	if err = a.pc.Touch(webapi.ApiCaptchaPage, qs); err != nil {
		return
	}
	// Fetch CAPTCHA code image
	qs = protocol.Params{}.WithNow("_t")
	session.CodeImage, err = a.pc.GetContent(webapi.ApiCaptchaCodeImage, qs)
	if err != nil {
		return
	}
	// Fetch CAPTCHA keys image
	session.KeysImage, err = a.pc.GetContent(webapi.ApiCaptchaAllKeyImage, qs)
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
	qs := protocol.Params{}.
		WithInt("id", index).
		WithNow("_t")
	return a.pc.GetContent(webapi.ApiCaptchaOneKeyImage, qs)
}

// CaptchaSubmit submits the CAPTCHA code to session.
func (a *Agent) CaptchaSubmit(session *CaptchaSession, code string) (err error) {
	// Get captcha sign
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().UnixNano())
	qs := protocol.Params{}.
		With("callback", cb).
		WithNow("_")
	resp := &webapi.CaptchaSignResponse{}
	if err = a.pc.CallJsonpApi(webapi.ApiCaptchaSign, qs, resp); err != nil {
		return
	}
	// Submit captcha code
	form := protocol.Params{}.
		With("ac", "security_code").
		With("type", "web").
		With("sign", resp.Sign).
		With("code", code).
		With("cb", session.callback).
		ToForm()
	return a.pc.CallJsonApi(webapi.ApiCaptchaSubmit, nil, form, &webapi.BasicResponse{})
}
