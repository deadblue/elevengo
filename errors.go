package elevengo

type implError struct {
	message string
}

func (e *implError) Error() string {
	return e.message
}

var (
	errUpstreamError = &implError{"call upstream API failed."}

	errFileCursorInvalid = &implError{"cursor must be created by FileCursor()"}
	errFileStatFailed    = &implError{"get file info failed"}

	errOfflineCursorInvalid = &implError{"cursor must be created by OfflineCursor()"}

	errCaptchaCodeIncorrect = &implError{"captcha code incorrect"}
)
