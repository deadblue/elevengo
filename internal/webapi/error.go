package webapi

import "errors"

var (
	ErrNotLogin = errors.New("user not login")

	ErrPasswordIncorrect    = errors.New("password incorrect")
	ErrLoginTwoStepVerify   = errors.New("requires two-step verification")
	ErrAccountNotBindMobile = errors.New("account not binds mobile")
	ErrCredentialExpired    = errors.New("credential expired")

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
		40101009: ErrPasswordIncorrect,
		40101010: ErrLoginTwoStepVerify,
		40101030: ErrAccountNotBindMobile,
		40101032: ErrCredentialExpired,
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
