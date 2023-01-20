package webapi

import "errors"

var (
	ErrNotLogin = errors.New("user not login")

	ErrOfflineInvalidLink = errors.New("invalid download link")
	ErrOfflineTaskExisted = errors.New("offline task existed")

	ErrOrderNotSupport = errors.New("file order not supported")

	ErrPasswordIncorrect    = errors.New("password incorrect")
	ErrLoginTwoStepVerify   = errors.New("requires two-step verification")
	ErrAccountNotBindMobile = errors.New("account not binds mobile")
	ErrCredentialInvalid    = errors.New("credential invalid")
	ErrSessionExited        = errors.New("session exited")

	ErrQrcodeExpired = errors.New("qrcode expired")

	// ErrUnexpected is the fall-back error whose code is not handled.
	ErrUnexpected = errors.New("unexpected error")

	// ErrExist means an item which you want to create is already existed.
	ErrExist = errors.New("target already exists")
	// ErrNotExist means an item which you find is not existed.
	ErrNotExist = errors.New("target does not exist")

	// ErrReachEnd means there are no more item.
	ErrReachEnd = errors.New("reach the end")

	ErrInvalidCursor = errors.New("invalid cursor")

	ErrUploadTooLarge = errors.New("upload reach the limit")

	ErrImportDirectory = errors.New("can not import directory")

	ErrDownloadEmpty = errors.New("can not get download URL")

	ErrDownloadDirectory = errors.New("can not download directory")

	ErrVideoNotReady = errors.New("video is not ready")

	ErrInvalidImportURI = errors.New("invalid import URI")

	errMap = map[int]error{
		// Normal errors
		99:     ErrNotLogin,
		990001: ErrNotLogin,
		// Offline errors
		10004: ErrOfflineInvalidLink,
		10008: ErrOfflineTaskExisted,
		// Dir errors
		20004: ErrExist,
		// Label errors
		21003: ErrExist,
		// File errors
		20130827: ErrOrderNotSupport,
		// Login errors
		40101009: ErrPasswordIncorrect,
		40101010: ErrLoginTwoStepVerify,
		40101030: ErrAccountNotBindMobile,
		40101032: ErrCredentialInvalid,
		40101037: ErrSessionExited,
		// QRCode errors
		40199002: ErrQrcodeExpired,
	}
)

func getError(code int) error {
	if err, found := errMap[code]; found {
		return err
	}
	return ErrUnexpected
}
