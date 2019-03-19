package elevengo

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

func (c *Client) setup() (err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	c.jar = jar
	c.client = &http.Client{
		Jar: jar,
	}
	c.info = &_UserInfo{}
	c.offline = &_OfflineToken{}
	return nil
}

func (c *Client) requestRaw(url string, qs *RequestParameters, body io.Reader) (data []byte, err error) {
	// append query string
	if qs != nil {
		index := strings.IndexRune(url, '?')
		if index < 0 {
			url = fmt.Sprintf("%s?%s", url, qs.QueryString())
		} else {
			url = fmt.Sprintf("%s&%s", url, qs.QueryString())
		}
	}
	// prepare request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return
	}
	if body == nil {
		req.Method = "GET"
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")

	// send request
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	// check response
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error :%d", resp.StatusCode)
	}

	r := resp.Body
	defer resp.Body.Close()
	enc := resp.Header.Get("Content-Encoding")
	if enc == "gzip" {
		r, _ = gzip.NewReader(r)
	} else if enc == "deflate" {
		r = flate.NewReader(r)
	}

	data, err = ioutil.ReadAll(r)
	return
}

func (c *Client) requestJson(url string, qs *RequestParameters, body io.Reader, result interface{}) (err error) {
	data, err := c.requestRaw(url, qs, body)
	if err != nil {
		return
	}
	if result == nil {
		return nil
	} else {
		return json.Unmarshal(data, result)
	}
}
