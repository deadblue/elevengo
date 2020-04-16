package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"math/rand"
	"time"
)

const (
	cookieUrl    = "https://115.com"
	cookieDomain = ".115.com"

	cookieUid  = "UID"
	cookieCid  = "CID"
	cookieSeid = "SEID"

	apiUserInfo = "https://my.115.com/"
)

/*
Credential contains required information that the upstream server uses to
authenticate a signed-in user.

In detail, three cookies are required: "UID", "CID", "SEID", you can find
them from your browser cookie storage.
*/
type Credential struct {
	UID  string
	CID  string
	SEID string
}

/*
Basic information of the signed-in user.
*/
type UserInfo struct {
	Id   int
	Name string
}

/*
Import credentials into agent.
*/
func (a *Agent) CredentialImport(cr *Credential) (err error) {
	cookies := map[string]string{
		cookieUid:  cr.UID,
		cookieCid:  cr.CID,
		cookieSeid: cr.SEID,
	}
	a.hc.SetCookies(cookieUrl, cookieDomain, cookies)
	return a.getUserInfo()
}

/*
Export credentials from agent, you can store it for future use.
*/
func (a *Agent) CredentialExport() (cr *Credential, err error) {
	if cookies := a.hc.Cookies(cookieUrl); cookies == nil || len(cookies) == 0 {
		err = errCredentialsNotExist
	} else {
		cr = &Credential{
			UID:  cookies[cookieUid],
			CID:  cookies[cookieCid],
			SEID: cookies[cookieSeid],
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
		a.ui = &UserInfo{}
	}
	a.ui.Id = result.Data.UserId
	a.ui.Name = result.Data.UserName
	return
}

/*
Get signed in user information, return nil if none signed in.
*/
func (a *Agent) User() (info *UserInfo) {
	if a.ui != nil {
		info = &UserInfo{
			Id:   a.ui.Id,
			Name: a.ui.Name,
		}
	}
	return
}
