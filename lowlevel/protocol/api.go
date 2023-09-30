package protocol

import (
	"bytes"
	"io"
	"time"

	"github.com/deadblue/elevengo/internal/util"
)

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

// ExecuteApi calls an API, and returns errors
func (c *Client) ExecuteApi(spec ApiSpec) (err error) {
	payload := spec.Payload()
	if spec.IsCrypto() {
		spec.SetCryptoKey(c.ecc.EncodeToken(time.Now().UnixMilli()))
		payload, _ = c.encryptPayload(payload)
	}
	// Perform HTTP request
	var body io.ReadCloser
	if body, err = c.internalCall(spec.Url(), payload); err != nil {
		return
	}
	defer util.QuietlyClose(body)
	// Handle response
	if spec.IsCrypto() {
		var data []byte
		if data, err = c.decryptBody(body); err != nil {
			return
		}
		return spec.Parse(bytes.NewReader(data))
	} else {
		return spec.Parse(body)
	}
}

func (c *Client) internalCall(url string, payload Payload) (body io.ReadCloser, err error) {
	c.v.Wait()
	defer c.v.ClockIn()
	// Prepare request
	if payload != nil {
		body, err = c.Post(url, payload)
	} else {
		body, err = c.Get(url, nil)
	}
	return
}
