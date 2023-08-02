package errors

var errorsMap = map[int]error{
	// Normal errors
	99:     ErrNotLogin,
	911:    ErrCaptchaRequired,
	990001: ErrNotLogin,
	// Offline errors
	10004: ErrOfflineInvalidLink,
	10008: ErrOfflineTaskExisted,
	// File errors
	20004: ErrExist,
	20022: ErrInvalidOperation,
	// Label errors
	21003: ErrExist,
	// Download errors
	50003: ErrNotExist,
	// Common errors
	990002: ErrInvalidParameters,
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

func Get(code int) error {
	if err, found := errorsMap[code]; found {
		return err
	}
	return ErrUnexpected
}
