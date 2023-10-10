package client

import (
	"bytes"
	"io"
)

const (
	_ContentTypeWwwForm = "application/x-www-form-urlencoded"
)

// Payload describes the request body.
type Payload interface {
	io.Reader

	// ContentType returns the MIME type of payload.
	ContentType() string

	// ContentLength returns the size in bytes of payload.
	ContentLength() int64
}

// PayloadImpl is a standard `Payload` implementation.
type _PayloadImpl struct {
	r io.Reader
	t string
	l int64
}

func (pi *_PayloadImpl) Read(p []byte) (int, error) {
	return pi.r.Read(p)
}

func (pi *_PayloadImpl) ContentType() string {
	return pi.t
}

func (pi *_PayloadImpl) ContentLength() int64 {
	return pi.l
}

func WwwFormPayload(s string) Payload {
	body := []byte(s)
	return &_PayloadImpl{
		r: bytes.NewReader(body),
		t: _ContentTypeWwwForm,
		l: int64(len(body)),
	}
}
