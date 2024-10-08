package elevengo

import (
	"errors"
	"iter"
)

var (
	errNoMoreItems = errors.New("no more items")
)

// Iterator iterate items.
type Iterator[T any] interface {

	// Count return the count of items.
	Count() int

	// Items return an index-item sequence.
	Items() iter.Seq2[int, T]
}
