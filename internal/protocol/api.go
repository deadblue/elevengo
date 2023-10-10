package protocol

import (
	"bytes"
	"io"
	"time"

	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/client"
)

func (c *ClientImpl) CallApi(spec client.ApiSpec) (err error) {
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

func (c *ClientImpl) internalCall(url string, payload client.Payload) (body io.ReadCloser, err error) {
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
