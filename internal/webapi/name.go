package webapi

import "fmt"

const (
	CookieDomain115   = ".115.com"
	CookieDomainAnxia = ".anxia.com"

	CookieUrl      = "https://115.com"
	CookieNameUid  = "UID"
	CookieNameCid  = "CID"
	CookieNameSeid = "SEID"

	defaultName = "Mozilla/5.0"
	appName     = "115Desktop"
)

var (
	AppVersion = ""
)

func MakeUserAgent(name string) string {
	if name == "" {
		name = defaultName
	}
	return fmt.Sprintf("%s %s/%s", name, appName, AppVersion)
}
