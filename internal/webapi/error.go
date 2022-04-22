package webapi

import "errors"

var (
	ErrNotLogin      = errors.New("user not login")
	ErrUnexpected    = errors.New("unexpected error")
	ErrVideoNotReady = errors.New("video is not ready")

	ErrOfflineInvalidLink = errors.New("invalid download link")
	ErrOfflineTaskExisted = errors.New("offline task existed")

	errMap = map[int]error{
		// Normal error
		990001: ErrNotLogin,
		// Offline error
		10004: ErrOfflineInvalidLink,
		10008: ErrOfflineTaskExisted,
	}
)

func getError(code int) error {
	if err, found := errMap[code]; found {
		return err
	}
	return ErrUnexpected
}
