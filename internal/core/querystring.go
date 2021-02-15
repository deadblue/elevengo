package core

import (
	"net/url"
	"strconv"
)

// A query-string holder
type QueryString interface {
	// Add a string value
	WithString(name, value string) QueryString
	// Add an int value
	WithInt(name string, value int) QueryString
	// Add an int64 value
	WithInt64(name string, value int64) QueryString
	// Add an uint64 value
	WithUint64(name string, value uint64) QueryString
	// Encode to query-string
	Encode() string
}

type implQueryString struct {
	values url.Values
}

func (i *implQueryString) WithString(name, value string) QueryString {
	i.values.Set(name, value)
	return i
}

func (i *implQueryString) WithInt(name string, value int) QueryString {
	return i.WithString(name, strconv.Itoa(value))
}

func (i *implQueryString) WithInt64(name string, value int64) QueryString {
	return i.WithString(name, strconv.FormatInt(value, 10))
}

func (i *implQueryString) WithUint64(name string, value uint64) QueryString {
	return i.WithString(name, strconv.FormatUint(value, 10))
}

func (i *implQueryString) Encode() string {
	return i.values.Encode()
}

func NewQueryString() QueryString {
	return &implQueryString{
		values: url.Values{},
	}
}
