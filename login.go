package elevengo

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
)

var (
	_RegexpUserId = regexp.MustCompile(`(?m)USER_ID = '([0-9]+)';`)
)

type Credentials struct {
	Uid  string
	Cid  string
	Seid string
}

func (c *Client) ImportCredentials(cr *Credentials) (err error) {
	cks := []*http.Cookie{
		{Name: "115_lang", Value: "zh", Domain: ".115.com", Path: "/", HttpOnly: false},
		{Name: "UID", Value: cr.Uid, Domain: ".115.com", Path: "/", HttpOnly: false},
		{Name: "CID", Value: cr.Cid, Domain: ".115.com", Path: "/", HttpOnly: false},
		{Name: "SEID", Value: cr.Seid, Domain: ".115.com", Path: "/", HttpOnly: false},
	}
	u, _ := url.Parse(apiHost)
	c.jar.SetCookies(u, cks)

	return c.check()
}

func (c *Client) check() (err error) {
	// request home page
	qs := newQueryString().WithString("mode", "wangpan")
	data, err := c.request(apiHost, qs, nil)
	if err != nil {
		return
	}
	// search and store user id
	body := string(data)
	matches := _RegexpUserId.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return errors.New("not login")
	} else {
		c.info.UserId = matches[0][1]
	}
	return nil
}
