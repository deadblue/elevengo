package elevengo

type ApiError struct {
	message string
}

func (e *ApiError) Error() string {
	return e.message
}

var (
	errUpstreamError = &ApiError{"call upstream API failed."}

	errInvalidFileCursor    = &ApiError{"cursor must be created by FileCursor()"}
	errInvalidOfflineCursor = &ApiError{"cursor must be created by OfflineCursor()"}
)
