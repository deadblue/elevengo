package protocol

import (
	"bytes"

	"github.com/deadblue/elevengo/lowlevel/client"
)

const (
	_ContentTypeWwwForm = "application/x-www-form-urlencoded"
)

// PayloadImpl is a standard `Payload` implementation.
type _PayloadImpl struct {
	r *bytes.Reader
	t string
}

func (pi *_PayloadImpl) Read(p []byte) (int, error) {
	return pi.r.Read(p)
}

func (pi *_PayloadImpl) ContentType() string {
	return pi.t
}

func (pi *_PayloadImpl) ContentLength() int64 {
	return pi.r.Size()
}

// WwwFormPayload constructs a www URL-encoded form payload.
func WwwFormPayload(s string) client.Payload {
	body := []byte(s)
	return &_PayloadImpl{
		r: bytes.NewReader(body),
		t: _ContentTypeWwwForm,
	}
}

// CustomPayload constructs payload with given data and MIME type.
func CustomPayload(data []byte, mimeType string) client.Payload {
	return &_PayloadImpl{
		r: bytes.NewReader(data),
		t: mimeType,
	}
}
