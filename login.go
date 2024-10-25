package elevengo

import (
	"context"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/api"
)

// Credential contains required information which 115 server uses to
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
	// User ID
	Id int
	// User name
	Name string
	// Is user VIP
	IsVip bool
}

// CredentialImport imports credentials into agent.
func (a *Agent) CredentialImport(cr *Credential) (err error) {
	cookies := map[string]string{
		protocol.CookieNameUID:  cr.UID,
		protocol.CookieNameCID:  cr.CID,
		protocol.CookieNameSEID: cr.SEID,
	}
	a.llc.ImportCookies(cookies, protocol.CookieDomains...)
	return a.afterSignIn(cr.UID)
}

// CredentialExport exports current credentials for future-use.
func (a *Agent) CredentialExport(cr *Credential) {
	cookies := a.llc.ExportCookies(protocol.CookieUrl)
	cr.UID = cookies[protocol.CookieNameUID]
	cr.CID = cookies[protocol.CookieNameCID]
	cr.SEID = cookies[protocol.CookieNameSEID]
}

func (a *Agent) afterSignIn(uid string) (err error) {
	// Call UploadInfo API to get userId and userKey
	spec := (&api.UploadInfoSpec{}).Init()
	if err = a.llc.CallApi(spec, context.Background()); err == nil {
		a.common.SetUserInfo(spec.Result.UserId, spec.Result.UserKey)
		a.isWeb = protocol.IsWebCredential(uid)
	}
	return
}

// UserGet get information of current signed-in user.
func (a *Agent) UserGet(info *UserInfo) (err error) {
	spec := (&api.UserInfoSpec{}).Init()
	if err = a.llc.CallApi(spec, context.Background()); err == nil {
		info.Id = spec.Result.UserId
		info.Name = spec.Result.UserName
		info.IsVip = spec.Result.IsVip != 0
	}
	return
}
