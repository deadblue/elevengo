package elevengo

type Iterator[T any] interface {
	Next() error
	Get(*T) error
}
