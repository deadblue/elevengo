package errors

import "errors"

var (
	ErrNotLogin = errors.New("user not login")

	ErrCaptchaRequired = errors.New("please resolve captcha")

	ErrOfflineInvalidLink = errors.New("invalid download link")

	ErrIPAbnormal           = errors.New("ip abnormal")
	ErrPasswordIncorrect    = errors.New("password incorrect")
	ErrLoginTwoStepVerify   = errors.New("requires two-step verification")
	ErrAccountNotBindMobile = errors.New("account not binds mobile")
	ErrCredentialInvalid    = errors.New("credential invalid")
	ErrSessionExited        = errors.New("session exited")

	ErrQrcodeExpired = errors.New("qrcode expired")
	ErrGetFailed     = errors.New("get failed")

	// ErrUnexpected is the fall-back error whose code is not handled.
	ErrUnexpected = errors.New("unexpected error")

	// ErrExist means an item which you want to create is already existed.
	ErrExist = errors.New("target already exists")
	// ErrNotExist means an item which you find is not existed.
	ErrNotExist = errors.New("target does not exist")

	ErrInvalidOperation = errors.New("invalid operation")

	ErrInvalidParameters = errors.New("invalid parameters")

	// ErrReachEnd means there are no more item.
	// ErrReachEnd = errors.New("reach the end")

	ErrUploadDisabled = errors.New("upload function is disabled")

	ErrUploadNothing = errors.New("nothing ot upload")

	ErrUploadTooLarge = errors.New("upload reach the limit")

	ErrInitUploadUnknowStatus = errors.New("unknown status from initupload")

	ErrImportDirectory = errors.New("can not import directory")

	ErrDownloadEmpty = errors.New("can not get download URL")

	ErrDownloadDirectory = errors.New("can not download directory")

	ErrVideoNotReady = errors.New("video is not ready")

	ErrEmptyList = errors.New("list is empty")
)
