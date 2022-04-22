package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
	"math/rand"
	"time"
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

func (u *UserInfo) IsLogin() bool {
	return u.Id != 0
}

// CredentialImport imports credentials into agent.
func (a *Agent) CredentialImport(cr *Credential) (err error) {
	cookies := map[string]string{
		webapi.CookieNameUid:  cr.UID,
		webapi.CookieNameCid:  cr.CID,
		webapi.CookieNameSeid: cr.SEID,
	}
	a.wc.ImportCookies(cookies, webapi.CookieDomain115, webapi.CookieDomainAnxia)
	return a.syncUserInfo()
}

// CredentialExport exports credentials for future-use.
func (a *Agent) CredentialExport(cr *Credential) (err error) {
	cookies := a.wc.ExportCookies(webapi.CookieUrl)
	cr.UID = cookies[webapi.CookieNameUid]
	cr.CID = cookies[webapi.CookieNameCid]
	cr.SEID = cookies[webapi.CookieNameSeid]
	return
}

// syncUserInfo syncs user information from cloud to agent.
func (a *Agent) syncUserInfo() (err error) {
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().Unix())
	qs := protocol.Params{}.
		With("callback", cb).
		WithNow("_")
	resp := webapi.BasicResponse{}
	if err = a.wc.CallJsonpApi(webapi.ApiUserInfo, qs, &resp); err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return err
	}
	result := webapi.UserInfoData{}
	if err = resp.Decode(&result); err != nil {
		return
	}
	a.user.Id = result.UserId
	a.user.Name = result.UserName
	return
}

// User returns user information.
func (a *Agent) User() *UserInfo {
	return &a.user
}
