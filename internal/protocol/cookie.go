package protocol

import "strings"

const (
	CookieUrl = "https://115.com"

	CookieNameUID  = "UID"
	CookieNameCID  = "CID"
	CookieNameSEID = "SEID"
)

var (
	CookieDomains = []string{
		".115.com",
		".anxia.com",
	}
)

func IsWebCredential(uid string) bool {
	parts := strings.Split(uid, "_")
	return len(parts) == 3 && parts[1][0] == 'A'
}
