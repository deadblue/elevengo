package types

import (
	"strconv"

	"github.com/deadblue/elevengo/internal/crypto/hash"
)

// Common parameters for several APIs
type CommonParams struct {
	// App version
	AppVer string
	// User ID
	UserId string
	// MD5 hash of user ID
	UserHash string
	// User key for uploading
	UserKey string
}

func (c *CommonParams) SetUserInfo(userId int, userKey string) {
	c.UserId = strconv.Itoa(userId)
	c.UserHash = hash.Md5Hex(c.UserId)
	c.UserKey = userKey
}
