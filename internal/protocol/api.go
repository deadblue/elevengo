package protocol

import (
	"bytes"
	"io"
	urllib "net/url"
	"strings"
	"time"

	"github.com/deadblue/elevengo/internal/util"
)

// ApiSpec describes the specification of an API.
type ApiSpec interface {
	IsCrypto() bool
	Url() string
	Payload() Payload
	Parse(r io.Reader) (err error)
}

func (c *Client) ExecuteApi(spec ApiSpec) (err error) {
	url, payload := spec.Url(), spec.Payload()
	if spec.IsCrypto() {
		// Append EC key to querystring
		ecKey := c.ecc.EncodeToken(time.Now().UnixMilli())
		if strings.ContainsRune(url, '?') {
			url = url + "&k_ec=" + urllib.QueryEscape(ecKey)
		} else {
			url = url + "?k_ec=" + urllib.QueryEscape(ecKey)
		}
		// Encrypt payload
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
			body, err = c.Post(url, nil, payload)
		} else {
			body, err = c.Get(url, nil, nil)
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
