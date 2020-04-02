package elevengo

type ApiError struct {
	message string
}

func (e *ApiError) Error() string {
	return e.message
}

var (
	ErrRemoteFailed = &ApiError{"call remote api failed"}
)
