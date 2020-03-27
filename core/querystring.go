package core

import (
	"net/url"
	"strconv"
	"time"
)

type QueryString struct {
	values url.Values
}

func NewQueryString() *QueryString {
	return &QueryString{
		values: url.Values{},
	}
}
func (qs *QueryString) WithString(name, value string) *QueryString {
	qs.values.Set(name, value)
	return qs
}
func (qs *QueryString) WithInt(name string, value int) *QueryString {
	return qs.WithString(name, strconv.Itoa(value))
}
func (qs *QueryString) WithInt64(name string, value int64) *QueryString {
	return qs.WithString(name, strconv.FormatInt(value, 10))
}
func (qs *QueryString) WithTimestamp(name string) *QueryString {
	return qs.WithInt64(name, time.Now().UnixNano())
}
func (qs *QueryString) Encode() string {
	return qs.values.Encode()
}
