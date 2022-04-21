package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"github.com/deadblue/elevengo/internal/webapi"
	"math/rand"
	"time"
)

const (
	apiUserInfo = "https://my.115.com/"
)

// Credential contains required information which upstream server uses to
// authenticate a signed-in user.
// In detail, three cookies are required: "UID", "CID", "SEID", caller can
// find them from browser cookie storage.
type Credential struct {
	UID  string
	CID  string
	SEID string
}

// UserInfo contains the basic information of a signed-in user.
type UserInfo struct {
	Id   int
	Name string
}

// CredentialImport imports credentials into agent.
func (a *Agent) CredentialImport(cr *Credential) (err error) {
	cookies := map[string]string{
		webapi.CookieNameUid:  cr.UID,
		webapi.CookieNameCid:  cr.CID,
		webapi.CookieNameSeid: cr.SEID,
	}
	a.pc.ImportCookies(cookies, webapi.CookieDomain115, webapi.CookieDomainAnxia)
	return nil
}

// CredentialExport exports credentials for future-use.
func (a *Agent) CredentialExport(cr *Credential) (err error) {
	cookies := a.pc.ExportCookies(webapi.CookieUrl)
	cr.UID = cookies[webapi.CookieNameUid]
	cr.CID = cookies[webapi.CookieNameCid]
	cr.SEID = cookies[webapi.CookieNameSeid]
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
	result := &types.UserInfoResult{}
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

// User returns user information.
func (a *Agent) User() (info UserInfo) {
	if a.ui != nil {
		info = *a.ui
	}
	return
}
