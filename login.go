package elevengo

import (
	"errors"
	"github.com/deadblue/elevengo/internal"
	"net/http"
	"net/url"
	"regexp"
)

const (
	cookieDomain = ".115.com"
	cookieUrl    = "https://115.com"

	urlUserInfo = "https://115.com/?mode=wangpan"
)

var (
	regexpUserId = regexp.MustCompile(`(?m)USER_ID = '([0-9]+)';`)
)

type Credentials struct {
	UID  string
	CID  string
	SEID string
}

func (c *Client) ImportCredentials(cr *Credentials) (err error) {
	cks := []*http.Cookie{
		{Name: "UID", Value: cr.UID, Domain: cookieDomain, Path: "/", HttpOnly: true},
		{Name: "CID", Value: cr.CID, Domain: cookieDomain, Path: "/", HttpOnly: true},
		{Name: "SEID", Value: cr.SEID, Domain: cookieDomain, Path: "/", HttpOnly: true},
	}
	u, _ := url.Parse(cookieUrl)
	c.cj.SetCookies(u, cks)
	return c.getUserData()
}

func (c *Client) getUserData() (err error) {
	// get home page
	data, err := c.hc.Get(urlUserInfo, nil)
	if err != nil {
		return
	}
	// search and store user id
	body := string(data)
	matches := regexpUserId.FindAllStringSubmatch(body, -1)
	if len(matches) == 0 {
		return errors.New("not login")
	}
	// store UserId
	if c.ui == nil {
		c.ui = new(internal.UserInfo)
	}
	c.ui.UserId = matches[0][1]
	return nil
}
