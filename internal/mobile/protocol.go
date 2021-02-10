package mobile

import (
	"fmt"
	"github.com/deadblue/gostream/quietly"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	urllib "net/url"
	"strconv"
)

var (
	_CookieUrl, _ = urllib.Parse("https://115.com")

	_UserAgent = fmt.Sprintf("Mozilla/5.0 (F5321; 8.0.0; en;) 115disk/%s", AppVersion)
)

func (c *Client) initHttpClient() {
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
	c.userId = uint32(userId)
}

func (c *Client) callApi(url string, apiId int) (err error) {
	params := urllib.Values{}
	params.Set("k_ec", c.ecEncodeKey(apiId))
	params.Set("user_id", strconv.FormatUint(uint64(c.userId), 10))
	params.Set("app_ver", AppVersion)

	req, _ := http.NewRequest(http.MethodGet,
		url+"?"+params.Encode(), nil)
	req.Header.Set("User-Agent", _UserAgent)

	resp, err := c.hc.Do(req)
	if err != nil {
		return
	}
	defer quietly.Close(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return c.ecDecode(body)
}
