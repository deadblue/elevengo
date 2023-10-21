package client

import "io"

// ApiSpec describes the specification of an 115 API.
type ApiSpec interface {

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
