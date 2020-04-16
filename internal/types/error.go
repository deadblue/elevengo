package types

import (
	"fmt"
)

type ErrorCategory string

const (
	FileError    ErrorCategory = "file"
	OfflineError ErrorCategory = "offline"
	QrcodeError  ErrorCategory = "qrcode"
)

// Common API error
type ApiError struct {
	Category ErrorCategory
	Code     int
	error    string
}

func (ae *ApiError) Error() string {
	return ae.error
}

func MakeFileError(code int, message string) *ApiError {
	return &ApiError{
		Category: FileError,
		Code:     code,
		error:    fmt.Sprintf("file API error[%d]: %s", code, message),
	}
}
func MakeOfflineError(code int, message string) *ApiError {
	return &ApiError{
		Category: OfflineError,
		Code:     code,
		error:    fmt.Sprintf("offline API error[%d]: %s", code, message),
	}
}
func MakeQrcodeError(code int, message string) *ApiError {
	return &ApiError{
		Category: QrcodeError,
		Code:     code,
		error:    fmt.Sprintf("qrcode API error[%d]: %s", code, message),
	}
}
