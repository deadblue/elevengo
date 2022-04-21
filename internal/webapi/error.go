package webapi

import "errors"

var (
	ErrNotLogin      = errors.New("user not login")
	ErrUnexpected    = errors.New("unexpected error")
	ErrVideoNotReady = errors.New("video is not ready")

	errMap = map[int]error{
		990001: ErrNotLogin,
	}
)

func getError(code int) error {
	if err, found := errMap[code]; found {
		return err
	}
	return ErrUnexpected
}
