package webapi

import "errors"

var (
	ErrNotLogin = errors.New("user not login")

	ErrLoginPasswordIncorrect = errors.New("password incorrect")
	ErrAccountNotBindMobile   = errors.New("account not binds mobile")

	ErrQrcodeExpired = errors.New("qrcode expired")

	ErrOfflineInvalidLink = errors.New("invalid download link")
	ErrOfflineTaskExisted = errors.New("offline task existed")

	// ErrUnexpected is the fall-back error whose code is not handled.
	ErrUnexpected = errors.New("unexpected error")

	ErrVideoNotReady = errors.New("video is not ready")

	errMap = map[int]error{
		// Normal errors
		990001: ErrNotLogin,
		// Offline errors
		10004: ErrOfflineInvalidLink,
		10008: ErrOfflineTaskExisted,
		// Login errors
		40101009: ErrLoginPasswordIncorrect,
		40101030: ErrAccountNotBindMobile,
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
