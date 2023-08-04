package errors

const (
	CodeFileOrderNotSupported = 20130827
)

type ErrFileOrderNotSupported struct {
	Order string
	Asc   int
}

func (e *ErrFileOrderNotSupported) Error() string {
	return "order not supported"
}
