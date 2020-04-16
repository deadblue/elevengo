package elevengo

import (
	"errors"
	"github.com/deadblue/elevengo/internal/types"
)

// Special errors
var (
	errCredentialsNotExist = errors.New("credentials not exist")

	errFileNotExist      = errors.New("file not exist")
	errFileCursorInvalid = errors.New("cursor must be created by FileCursor()")

	errOfflineCursorInvalid = errors.New("cursor must be created by OfflineCursor()")

	errCaptchaFailed = errors.New("captcha code incorrect")

	errVideoNotReady = errors.New("video not ready")
)

const (
	codeFileExist = 20004

	codeOfflineCaptcha = 911
	codeOfflineExist   = 10008

	codeQrcodeExpired = 40199002
)

func IsFileExist(err error) bool {
	if ae, ok := err.(*types.ApiError); ok {
		return ae.Category == types.FileError && ae.Code == codeFileExist
	}
	return false
}

func IsFileNotExist(err error) bool {
	return err == errFileNotExist
}

func IsOfflineExist(err error) bool {
	if ae, ok := err.(*types.ApiError); ok {
		return ae.Category == types.OfflineError && ae.Code == codeOfflineExist
	}
	return false
}

func IsOfflineCaptcha(err error) bool {
	if ae, ok := err.(*types.ApiError); ok {
		return ae.Category == types.OfflineError && ae.Code == codeOfflineCaptcha
	}
	return false
}

func IsQrcodeExpire(err error) bool {
	if ae, ok := err.(*types.ApiError); ok {
		return ae.Category == types.QrcodeError && ae.Code == codeQrcodeExpired
	}
	return false
}
