package base

import (
	"bytes"
	"io"

	"github.com/deadblue/elevengo/internal/protocol"
)

const (
	contentTypeWwwForm = "application/x-www-form-urlencoded"
)

type PayloadImpl struct {
	r io.Reader
	t string
	s int64
}

func (pi *PayloadImpl) Read(p []byte) (int, error) {
	return pi.r.Read(p)
}

func (pi *PayloadImpl) ContentType() string {
	return pi.t
}

func (pi *PayloadImpl) ContentLength() int64 {
	return pi.s
}

func wwwFormPayload(s string) protocol.Payload {
	body := []byte(s)
	return &PayloadImpl{
		r: bytes.NewReader(body),
		t: contentTypeWwwForm,
		s: int64(len(body)),
	}
}
