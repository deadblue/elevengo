package elevengo

import (
	"errors"
	"github.com/deadblue/elevengo/internal"
)

// Special errors
var (
	errFileNotExist      = errors.New("file not exist")
	errFileCursorInvalid = errors.New("cursor must be created by FileCursor()")

	errOfflineCursorInvalid = errors.New("cursor must be created by OfflineCursor()")

	errCaptchaFailed = errors.New("captcha code incorrect")
)

const (
	codeFileExist = 20004

	codeOfflineCaptcha = 911
	codeOfflineExist   = 10008

	codeQrcodeExpired = 40199002
)

func IsFileExist(err error) bool {
	if ae, ok := err.(*internal.ApiError); ok {
		return ae.Category == internal.FileError && ae.Code == codeFileExist
	}
	return false
}

func IsFileNotExist(err error) bool {
	return err == errFileNotExist
}

func IsOfflineExist(err error) bool {
	if ae, ok := err.(*internal.ApiError); ok {
		return ae.Category == internal.OfflineError && ae.Code == codeFileExist
	}
	return false
}

func IsOfflineCaptcha(err error) bool {
	if ae, ok := err.(*internal.ApiError); ok {
		return ae.Category == internal.OfflineError && ae.Code == codeOfflineCaptcha
	}
	return false
}

func IsQrcodeExpire(err error) bool {
	if ae, ok := err.(*internal.ApiError); ok {
		return ae.Category == internal.QrcodeError && ae.Code == codeQrcodeExpired
	}
	return false
}
