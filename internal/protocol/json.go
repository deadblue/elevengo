package protocol

import (
	"encoding/json"
	"github.com/deadblue/gostream/quietly"
	"io"
)

// CallJsonApi calls remote HTTP API, and parses its result as JSON.
func (c *Client) CallJsonApi(url string, qs Params, form Params, result interface{}) (err error) {
	// Prepare request
	var body io.ReadCloser
	if form != nil {
		body, err = c.PostForm(url, qs, form)
	} else {
		body, err = c.Get(url, qs)
	}
	if err != nil {
		return
	}
	defer quietly.Close(body)
	// Parse response
	if result != nil {
		decoder := json.NewDecoder(body)
		err = decoder.Decode(result)
	}
	return
}

func (c *Client) JsonPApi(url string, qs Params, result interface{}) (err error) {
	return
}
