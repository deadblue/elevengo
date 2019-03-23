package elevengo

import (
	"fmt"
)

func apiError(code int) error {
	return fmt.Errorf("api error: %d", code)
}

type ElevenFiveError struct {
	Message string
}

func (e *ElevenFiveError) Error() string {
	return e.Message
}

var (
	ErrInvalidResult = &ElevenFiveError{"invalid API result"}
	ErrCaptchaFailed = &ElevenFiveError{"captcha code incorrect"}
)
