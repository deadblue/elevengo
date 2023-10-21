package client

import (
	"io"
)

// Payload describes the request body.
type Payload interface {
	io.Reader

	// ContentType returns the MIME type of payload.
	ContentType() string

	// ContentLength returns the size in bytes of payload.
	ContentLength() int64
}
