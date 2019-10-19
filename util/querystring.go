package util

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
func (qs *QueryString) WithTimestamp(name string) *QueryString {
	value := strconv.FormatInt(time.Now().UnixNano(), 10)
	return qs.WithString(name, value)
}
func (qs *QueryString) Encode() string {
	return qs.values.Encode()
}
