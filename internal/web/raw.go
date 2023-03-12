package web

import (
	"bytes"
)

type rawPayload struct {
	t string
	r *bytes.Reader
}

func (r *rawPayload) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (r *rawPayload) ContentType() string {
	return r.t
}

func (r *rawPayload) ContentLength() int64 {
	return r.r.Size()
}

func makePayload(data []byte, contentType string) Payload {
	return &rawPayload{
		t: contentType,
		r: bytes.NewReader(data),
	}
}