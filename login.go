package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"math/rand"
	"time"
)

const (
	cookieDomain = ".115.com"
	cookieUrl    = "https://115.com"

	apiUserInfo = "https://my.115.com/"
)

/*
Credentials contains required information that the upstream server uses to
authenticate a signed-in user.

In detail, three cookies are required: "UID", "CID", "SEID", you can find
them from your browser cookie storage.
*/
type Credentials struct {
	UID  string
	CID  string
	SEID string
}

/*
Import credentials into agent.
*/
func (a *Agent) CredentialsImport(cr *Credentials) (err error) {
	cookies := map[string]string{
		"UID":  cr.UID,
		"CID":  cr.CID,
		"SEID": cr.SEID,
	}
	a.hc.SetCookies(cookieUrl, cookieDomain, cookies)
	return a.getUserInfo()
}

/*
Export credentials from agent, you can store it for future use.
*/
func (a *Agent) CredentialsExport() (cr *Credentials, err error) {
	if cookies := a.hc.Cookies(cookieUrl); cookies == nil || len(cookies) == 0 {
		err = errCredentialsNotExist
	} else {
		cr = &Credentials{
			UID:  cookies["UID"],
			CID:  cookies["CID"],
			SEID: cookies["SEID"],
		}
	}
	return
}

// A new and graceful way to get user information.
func (a *Agent) getUserInfo() (err error) {
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().Unix())
	qs := core.NewQueryString().
		WithString("ct", "ajax").
		WithString("ac", "nav").
		WithString("callback", cb).
		WithInt64("_", time.Now().Unix())
	result := &internal.UserInfoResult{}
	if err = a.hc.JsonpApi(apiUserInfo, qs, result); err != nil {
		return
	}
	if a.ui == nil {
		a.ui = new(internal.UserInfo)
	}
	a.ui.UserId = result.Data.UserId
	a.ui.UserName = result.Data.UserName
	return
}
