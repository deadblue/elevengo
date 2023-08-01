package elevengo

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
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
	return a.afterSignIn()
}

// CredentialExport exports credentials for future-use.
func (a *Agent) CredentialExport(cr *Credential) {
	cookies := a.pc.ExportCookies(webapi.CookieUrl)
	cr.UID = cookies[webapi.CookieNameUid]
	cr.CID = cookies[webapi.CookieNameCid]
	cr.SEID = cookies[webapi.CookieNameSeid]
}

func (a *Agent) afterSignIn() (err error) {
	// Call UploadInfo API to get userId and userKey
	uis := (&api.UploadInfoSpec{}).Init()
	if err = a.pc.ExecuteApi(uis); err != nil {
		return
	} else {
		a.uh.SetUserData(uis.Resp.UserId, uis.Resp.UserKey)
	}
	// Call IndexInfo to get session information
	iis := (&api.IndexInfoSpec{}).Init()
	if err = a.pc.ExecuteApi(uis); err != nil {
		return
	}
	for _, li := range iis.Resp.Data.LoginInfos.List {
		if li.IsCurrent == 1 {
			a.isWeb = strings.HasPrefix(li.AppFlag, "A")
			break
		}
	}
	return
}

// UserGet retrieves user information from cloud.
func (a *Agent) UserGet(info *UserInfo) (err error) {
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().Unix())
	qs := protocol.Params{}.
		With("callback", cb).
		WithNow("_")
	resp := webapi.BasicResponse{}
	if err = a.pc.CallJsonpApi(webapi.ApiUserInfo, qs, &resp); err != nil {
		return
	}
	result := webapi.UserInfoData{}
	if err = resp.Decode(&result); err == nil {
		info.Id = result.UserId
		info.Name = result.UserName
	}
	return
}
