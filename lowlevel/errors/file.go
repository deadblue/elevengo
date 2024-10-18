package errors

const (
	CodeFileOrderNotSupported = 20130827
)

type FileOrderInvalidError struct {
	Order string
	Asc   int
}

func (e *FileOrderInvalidError) Error() string {
	return "invalid file order"
}
