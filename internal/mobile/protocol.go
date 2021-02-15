package mobile

import (
	"encoding/json"
	"fmt"
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/gostream/quietly"
	"io"
	"net/http"
	"net/http/cookiejar"
	urllib "net/url"
	"strings"
)

const (
	AppVersion = "26.1.0"
)

var (
	_CookieUrl, _ = urllib.Parse("https://115.com")

	_UserAgent = fmt.Sprintf("Mozilla/5.0 (F5321; 8.0.0; en;) 115disk/%s", AppVersion)
)

type _MobileResponse struct {
	State   bool            `json:"state"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (c *Client) httpInit() {
	c.cj, _ = cookiejar.New(nil)
	// TODO: Adjust the parameters of http client.
	c.hc = &http.Client{
		Jar: c.cj,
	}
}

/*
ImportCredentials imports mobile credentials (aka. cookies) into client.
Caller must use the cookie dumped from mobile client, not from browser.
*/
func (c *Client) ImportCredentials(userId int, credentials map[string]string) {
	c.uid = uint32(userId)
	cookies := make([]*http.Cookie, 0)
	for name, value := range credentials {
		cookies = append(cookies, &http.Cookie{
			Name:   name,
			Value:  value,
			Domain: ".115.com",
			Path:   "/",
			Secure: false,
		})
	}
	c.cj.SetCookies(_CookieUrl, cookies)
}

func (c *Client) callApi(url string, params map[string]string, data []byte, result interface{}) (err error) {
	// Make full URL
	qs := core.NewQueryString().
		WithString("k_ec", c.ecEncodeKey(0)).
		WithUint64("user_id", uint64(c.uid)).
		WithString("app_ver", AppVersion)
	if params != nil {
		for key, val := range params {
			qs.WithString(key, val)
		}
	}
	if strings.IndexRune(url, '?') < 0 {
		url = fmt.Sprintf("%s?%s", url, qs.Encode())
	} else {
		url = fmt.Sprintf("%s&%s", url, qs.Encode())
	}
	// Make request
	method, body := http.MethodGet, io.Reader(nil)
	if data != nil {
		method, body = http.MethodPost, c.ecEncode(data)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", _UserAgent)
	// Send request
	resp, err := c.hc.Do(req)
	if err != nil {
		return
	}
	defer quietly.Close(resp.Body)
	// Decrypt and parse response
	if result != nil {
		err = c.ecDecode(resp.Body, result)
	}
	return
}
