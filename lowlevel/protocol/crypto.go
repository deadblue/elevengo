package protocol

import (
	"bytes"
	"io"
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

func (c *Client) encryptPayload(p Payload) (ep Payload, err error) {
	if p == nil {
		return nil, nil
	}
	// Read payload
	body, err := io.ReadAll(p)
	if err != nil {
		return
	}
	// Encrypt it
	body = c.ecc.Encode(body)
	ep = &rawPayload{
		t: p.ContentType(),
		r: bytes.NewReader(body),
	}
	return
}

func (c *Client) decryptBody(r io.Reader) (data []byte, err error) {
	data, err = io.ReadAll(r)
	if err != nil {
		return
	}
	return c.ecc.Decode(data)
}
