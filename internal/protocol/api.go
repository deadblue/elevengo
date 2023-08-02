package protocol

import (
	"bytes"
	"io"
	"time"

	"github.com/deadblue/elevengo/internal/util"
)

// ApiSpec describes the specification of an API.
type ApiSpec interface {
	IsCrypto() bool
	SetCryptoKey(key string)
	Url() string
	Payload() Payload
	Parse(r io.Reader) (err error)
}

func (c *Client) ExecuteApi(spec ApiSpec) (err error) {
	payload := spec.Payload()
	if spec.IsCrypto() {
		spec.SetCryptoKey(c.ecc.EncodeToken(time.Now().UnixMilli()))
		payload, _ = c.encryptPayload(payload)
	}
	// Perform HTTP request
	var body io.ReadCloser
	{
		// Request frequency control
		c.v.Wait()
		defer c.v.ClockIn()
		// Prepare request
		if payload != nil {
			body, err = c.Post(spec.Url(), nil, payload)
		} else {
			body, err = c.Get(spec.Url(), nil, nil)
		}
	}
	if err != nil {
		return
	}
	// Handle response
	defer util.QuietlyClose(body)
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
