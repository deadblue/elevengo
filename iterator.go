package elevengo

import "github.com/deadblue/elevengo/internal/webapi"

// Iterator iterate items.
type Iterator[T any] interface {

	// Next move cursor to next.
	Next() error

	// Get gets current item.
	Get(*T) error

	// Count return the count of items.
	Count() int
}

func IsIteratorEnd(err error) bool {
	return err == webapi.ErrReachEnd
}
