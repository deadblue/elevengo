package elevengo

type implError struct {
	message string
}

func (e *implError) Error() string {
	return e.message
}

var (
	errUpstreamError = &implError{"call upstream API failed."}

	errInvalidFileCursor    = &implError{"cursor must be created by FileCursor()"}
	errInvalidOfflineCursor = &implError{"cursor must be created by OfflineCursor()"}

	errCaptchaCodeIncorrect = &implError{"captcha code incorrect"}
)
