package errors

import "fmt"

var errorsMap = map[int]error{
	// Normal errors
	99:     ErrNotLogin,
	911:    ErrCaptchaRequired,
	990001: ErrNotLogin,
	// Offline errors
	CodeOfflineIllegalLink: ErrOfflineInvalidLink,
	// File errors
	20004: ErrExist,
	20022: ErrInvalidOperation,
	// Label errors
	21003: ErrExist,
	// Download errors
	50003: ErrNotExist,
	// Common errors
	990002: ErrInvalidParameters,
	// Login errors
	40101004: ErrIPAbnormal,
	40101009: ErrPasswordIncorrect,
	40101010: ErrLoginTwoStepVerify,
	40101030: ErrAccountNotBindMobile,
	40101032: ErrCredentialInvalid,
	40101037: ErrSessionExited,
	// QRCode errors
	40199002: ErrQrcodeExpired,
	50199004: ErrGetFailed,

	// Whitelist errors
	CodeOfflineTaskExists: nil,
}

type ApiError struct {
	Code    int
	Message string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("(%d)%s", e.Code, e.Message)
}

func Get(code int, message string) error {
	if err, found := errorsMap[code]; found {
		return err
	}
	return &ApiError{
		Code:    code,
		Message: message,
	}
}
