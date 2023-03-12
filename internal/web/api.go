package web

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/deadblue/elevengo/internal/util"
)

type ApiResp interface {
	Err() error
}

// CallJsonApi calls remote HTTP API, and parses its result as JSON.
func (c *Client) CallJsonApi(
	url string, qs Params, payload Payload, 
	resp ApiResp,
) (err error) {
	// Prepare request
	var body io.ReadCloser
	if payload != nil {
		body, err = c.Post(url, qs, payload)
	} else {
		body, err = c.Get(url, qs, nil)
	}
	if err != nil {
		return
	}
	defer util.QuietlyClose(body)
	// Parse response
	if resp == nil {
		return
	}
	decoder := json.NewDecoder(body)
	if err = decoder.Decode(resp); err == nil {
		err = resp.Err()
	}
	return
}

func (c *Client) CallJsonpApi(url string, qs Params, resp ApiResp) (err error) {
	body, err := c.Get(url, qs, nil)
	if err != nil {
		return
	}
	defer util.QuietlyClose(body)
	data, err := io.ReadAll(body)
	if err != nil {
		return
	}
	left, right := bytes.IndexByte(data, '('), bytes.LastIndexByte(data, ')')
	if left < 0 || right < 0 {
		return &json.SyntaxError{Offset: 0}
	}
	if err = json.Unmarshal(data[left+1:right], resp); err == nil {
		err = resp.Err()
	}
	return
}

// CallSecretJsonApi calls JSON API with EC cryptography.
func (c *Client) CallSecretJsonApi(
	url string, qs Params, payload Payload, 
	resp ApiResp, timestamp int64,
) (err error) {
	var data []byte
	var body io.ReadCloser
	// Append EC key in querystring
	if qs == nil {
		qs = Params{}
	}
	qs.With("k_ec", c.ecc.EncodeToken(timestamp))
	// Encrypt payload
	if data, err = io.ReadAll(payload); err != nil {
		return
	}
	payload = makePayload(c.ecc.Encode(data), payload.ContentType())
	// Call API
	if body, err = c.Post(url, qs, payload); err != nil {
		return
	}
	defer util.QuietlyClose(body)
	if resp == nil {
		return
	}
	// Decrypt body
	if data, err = io.ReadAll(body); err != nil {
		return
	}
	if data, err = c.ecc.Decode(data); err != nil {
		return
	}
	// Parse result
	if err = json.Unmarshal(data, resp); err == nil {
		err = resp.Err()
	}
	return
}