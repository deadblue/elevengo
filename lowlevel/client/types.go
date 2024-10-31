package client

import "io"

type (
	// Payload describes the request body.
	Payload interface {
		io.Reader

		// ContentType returns the MIME type of payload.
		ContentType() string

		// ContentLength returns the size in bytes of payload.
		ContentLength() int64
	}

	// ApiSpec describes the specification of an 115 API.
	ApiSpec interface {

		// IsCrypto indicates whether the API request uses EC-crypto.
		IsCrypto() bool

		// SetCryptoKey adds crypto key in parameters.
		SetCryptoKey(key string)

		// Url returns the request URL of API.
		Url() string

		// Payload returns the request body of API.
		Payload() Payload

		// Parse parses the response body.
		Parse(r io.Reader) (err error)
	}

	Body interface {
		io.ReadCloser

		// Size returns body size or -1 when unknown.
		Size() int64

		// TotalSize returns total size of remote content or -1 when unknown.
		TotalSize() int64
	}
)
