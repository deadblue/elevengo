package elevengo

import (
	"encoding/base64"
	"fmt"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/elevengo/internal/webapi/sso"
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

// CredentialImport imports credentials into agent.
func (a *Agent) CredentialImport(cr *Credential) (err error) {
	cookies := map[string]string{
		webapi.CookieNameUid:  cr.UID,
		webapi.CookieNameCid:  cr.CID,
		webapi.CookieNameSeid: cr.SEID,
	}
	a.wc.ImportCookies(cookies, webapi.CookieDomain115, webapi.CookieDomainAnxia)
	if !a.LoginCheck() {
		err = webapi.ErrCredentialInvalid
	}
	return
}

// CredentialExport exports credentials for future-use.
func (a *Agent) CredentialExport(cr *Credential) {
	cookies := a.wc.ExportCookies(webapi.CookieUrl)
	cr.UID = cookies[webapi.CookieNameUid]
	cr.CID = cookies[webapi.CookieNameCid]
	cr.SEID = cookies[webapi.CookieNameSeid]
}

func (a *Agent) LoginCheck() bool {
	qs := web.Params{}.WithNowMilli("_")
	resp := &webapi.LoginBasicResponse{}
	if err := a.wc.CallJsonApi(webapi.ApiLoginCheck, qs, nil, resp); err != nil {
		return false
	}
	if resp.State != 0 {
		return false
	}
	data := &webapi.LoginCheckData{}
	if err := resp.Decode(data); err != nil {
		return false
	}
	a.uid = data.UserId
	return true
}

// UserGet retrieves user information from cloud.
func (a *Agent) UserGet(info *UserInfo) (err error) {
	cb := fmt.Sprintf("jQuery%d_%d", rand.Uint64(), time.Now().Unix())
	qs := web.Params{}.
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
	if err = resp.Decode(&result); err == nil {
		info.Id = result.UserId
		info.Name = result.UserName
	}
	return
}

func (a *Agent) loginGetKey() (key string, err error) {
	resp := &webapi.LoginBasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiLoginGetKey, nil, nil, resp); err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	data := &webapi.LoginKeyData{}
	if err = resp.Decode(data); err != nil {
		return
	}
	var keyData []byte
	if keyData, err = base64.StdEncoding.DecodeString(data.Key); err == nil {
		key = string(keyData)
	}
	return
}

// Login authenticates user by account and password.
//
// This API does not work for VIP user, because VIP user always need verify by SMS.
func (a *Agent) Login(account, password string) (err error) {
	// Get Login key
	key, err := a.loginGetKey()
	if err != nil {
		return
	}
	// Encrypt password
	now := time.Now().Unix()
	encPwd, err := sso.EncryptPassword(password, now, key)
	if err != nil {
		return
	}
	// Send login request
	ext := sso.GenerateExt()
	form := web.Params{}.
		With("login[version]", "2.0").
		With("login[safe]", "1").
		With("login[time]", "0").
		With("login[safe_login]", "0").
		With("login[country]", "").
		With("login[ssoent]", "A1").
		With("login[ssoext]", ext).
		With("login[ssovcode]", ext).
		With("login[ssoln]", account).
		With("login[ssopw]", sso.EncodePassword(account, password, ext)).
		WithInt("login[pwd_level]", sso.GetPasswordLevel(password)).
		With("goto", "https://115.com").
		With("country", "").
		With("from_browser", "1").
		With("cipher_ver", "2").
		With("account", account).
		With("passwd", encPwd).
		WithInt64("time", now)
	resp := &webapi.LoginBasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiPasswordLogin, nil, form, resp); err != nil {
		return
	}
	// Parse response
	if err = resp.Err(); err != nil {
		return
	}
	data := &webapi.LoginUserData{}
	if err = resp.Decode(data); err == nil {
		a.uid = data.Id
	}
	return
}

func (a *Agent) LoginSendSmsCode(userId int) (err error) {
	form := web.Params{}.
		With("tpl", "login_from_two_step").
		With("cv21", "2").
		WithInt("user_id", userId)
	resp := &webapi.LoginBasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiSmsSendCode, nil, form, resp); err == nil {
		err = resp.Err()
	}
	return
}

func (a *Agent) LoginBySms(userId int, code string) (err error) {
	form := web.Params{}.
		WithInt("account", userId).
		With("code", code)
	resp := &webapi.LoginBasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiSmsLogin, nil, form, resp); err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	data := &webapi.LoginUserData{}
	if err = resp.Decode(data); err == nil {
		a.uid = data.Id
	}
	return
}
