package elevengo

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func (c *Client) request(url string, qs *_QueryString, form *_Form) (data []byte, err error) {
	// append query string
	if qs != nil {
		index := strings.IndexRune(url, '?')
		if index < 0 {
			url = fmt.Sprintf("%s?%s", url, qs.Encode())
		} else {
			url = fmt.Sprintf("%s&%s", url, qs.Encode())
		}
	}
	// make request
	method, body := "", io.Reader(nil)
	if form == nil {
		method, body = http.MethodGet, nil
	} else {
		method, body = http.MethodPost, form.Finish()
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	// set request headers
	if form != nil {
		req.Header.Set("Content-Type", form.ContentType())
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", apiBasic)
	req.Header.Set("User-Agent", c.userAgent)
	// send request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	// check response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error :%d", resp.StatusCode)
	}
	// read body
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return
}

func (c *Client) requestJson(url string, qs *_QueryString, form *_Form, result interface{}) (err error) {
	data, err := c.request(url, qs, form)
	if err != nil {
		return
	}
	if result == nil {
		return nil
	} else {
		return json.Unmarshal(data, result)
	}
}
