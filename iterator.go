package elevengo

import "github.com/deadblue/elevengo/lowlevel/errors"

// Iterator iterate items.
type Iterator[T any] interface {

	// Next move cursor to next.
	Next() error

	// Index returns the index of current item, starts from 0.
	Index() int

	// Get gets current item.
	Get(*T) error

	// Count return the count of items.
	Count() int
}

func IsIteratorEnd(err error) bool {
	return err == errors.ErrReachEnd
}
