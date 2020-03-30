package elevengo

//import (
//	"fmt"
//)
//
//func apiError(code int) error {
//	return fmt.Errorf("api error: %d", code)
//}
//
//type ElevenFiveError struct {
//	Message string
//}
//
//func (e *ElevenFiveError) Error() string {
//	return e.Message
//}
//
//var (
//	ErrUnexpected          = &ElevenFiveError{"unexpected"}
//	ErrInvalidResult       = &ElevenFiveError{"invalid API result"}
//	ErrEmptyKeyword        = &ElevenFiveError{"empty key word"}
//	ErrUploadDirectory     = &ElevenFiveError{"can not upload directory"}
//	ErrCaptchaFailed       = &ElevenFiveError{"captcha code incorrect"}
//	ErrOfflineNothindToAdd = &ElevenFiveError{"nothing files to download!"}
//)
