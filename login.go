package elevengo

import (
	"strings"

	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/api/errors"
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
		api.CookieNameUID:  cr.UID,
		api.CookieNameCID:  cr.CID,
		api.CookieNameSEID: cr.SEID,
	}
	a.pc.ImportCookies(cookies, api.CookieDomains...)
	return a.afterSignIn(cr.UID)
}

// CredentialExport exports credentials for future-use.
func (a *Agent) CredentialExport(cr *Credential) {
	cookies := a.pc.ExportCookies(api.CookieUrl)
	cr.UID = cookies[api.CookieNameUID]
	cr.CID = cookies[api.CookieNameCID]
	cr.SEID = cookies[api.CookieNameSEID]
}

func (a *Agent) afterSignIn(uid string) (err error) {
	// Call UploadInfo API to get userId and userKey
	uis := (&api.UploadInfoSpec{}).Init()
	if err = a.pc.ExecuteApi(uis); err != nil {
		return
	} else {
		a.uh.SetUserParams(uis.Result.UserId, uis.Result.UserKey)
	}
	// Check UID
	parts := strings.Split(uid, "_")
	if len(parts) != 3 {
		return errors.ErrCredentialInvalid
	}
	a.isWeb = strings.HasPrefix(parts[1], "A")
	return
}

// UserGet retrieves user information from cloud.
// func (a *Agent) UserGet(info *UserInfo) (err error) {
// 	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().Unix())
// 	qs := protocol.Params{}.
// 		With("callback", cb).
// 		WithNow("_")
// 	resp := webapi.BasicResponse{}
// 	if err = a.pc.CallJsonpApi(webapi.ApiUserInfo, qs, &resp); err != nil {
// 		return
// 	}
// 	result := webapi.UserInfoData{}
// 	if err = resp.Decode(&result); err == nil {
// 		info.Id = result.UserId
// 		info.Name = result.UserName
// 	}
// 	return
// }
