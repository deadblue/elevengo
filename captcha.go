package elevengo

import (
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

// CaptchaSession holds the information during a CAPTCHA process.
type CaptchaSession struct {
	/*
		The CAPTCHA image data.

		There are 4 Chinese characters on this image, you need call "CaptchaKeysImage"
		to get a image which consists of 10 Chinese characters, then find 4 characters
		from them which matches with this image, the indexes of the 4 characters is
		the CAPTCHA code. (Index bases on zero.)
	*/
	CodeImage []byte

	/*
		TODO: Add doc.
	*/
	KeysImage []byte

	// The callback function name, hide for caller.
	callback string
}

/*
Start a CAPTCHA session.

Example:

	// Start CAPTCHA session
	session, err := agent.CaptchaStart()
	if err != nil {
		panic(err)
	}
	keysImg, err := agent.CaptchaKeysImage(session)
	if err != nil {
		panic(err)
	}
	// TODO: Solve the CAPTCHA here
	// Submit CAPTCHA code
	if err = agent.CaptchaSubmit(session, code); err != nil {
		panic(err)
	}

*/
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
Get CAPTCHA keys image data.

There are 10 Chinese characters on the image, 5 in column and 2 in row.
You can call this method multiple times, it will return the same 10
characters in different font on every calling.
*/
//func (a *Agent) CaptchaKeysImage(session *CaptchaSession) ([]byte, error) {
//	qs := core.NewQueryString().
//		WithString("ct", "index").
//		WithString("ac", "code").
//		WithString("t", "all").
//		WithInt64("_t", time.Now().Unix())
//	return a.hc.Get(apiCaptcha, qs)
//}

/*
Get one CAPTCHA key image data.

You can call this method multiple times, it will return the same character
in different font on every calling.
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
	signResult := &internal.CaptchaSignResult{}
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
	submitResult := &internal.CaptchaSubmitResult{}
	err = a.hc.JsonApi(apiCaptchaSubmit, nil, form, submitResult)
	if err == nil && submitResult.IsFailed() {
		err = errCaptchaCodeIncorrect
	}
	return
}
