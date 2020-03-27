package elevengo

import (
	"errors"
	"github.com/deadblue/elevengo/core"
	"net/http"
	"net/url"
	"regexp"
)

var (
	_RegexpUserId = regexp.MustCompile(`(?m)USER_ID = '([0-9]+)';`)
)

func (c *Client) ImportCredentials(cr *Credentials) (err error) {
	cks := []*http.Cookie{
		{Name: "115_lang", Value: "zh", Domain: domain, Path: "/", HttpOnly: false},
		{Name: "UID", Value: cr.UID, Domain: domain, Path: "/", HttpOnly: false},
		{Name: "CID", Value: cr.CID, Domain: domain, Path: "/", HttpOnly: false},
		{Name: "SEID", Value: cr.SEID, Domain: domain, Path: "/", HttpOnly: false},
	}
	u, _ := url.Parse(apiBasic)
	c.jar.SetCookies(u, cks)

	return c.getUserData()
}

func (c *Client) getUserData() (err error) {
	// request home page
	qs := core.NewQueryString().WithString("mode", "wangpan")
	data, err := c.request(apiBasic, qs, nil)
	if err != nil {
		return
	}
	// search and store user id
	body := string(data)
	matches := _RegexpUserId.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return errors.New("not login")
	}
	// store UserId
	if c.info == nil {
		c.info = &_UserInfo{}
	}
	c.info.UserId = matches[0][1]
	return nil
}
