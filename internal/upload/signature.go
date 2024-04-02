package upload

import (
	"crypto/sha1"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
)

func CalcSignature(userId, userKey, fileId, target string) string {
	digester := sha1.New()
	wx := util.UpgradeWriter(digester)
	// First pass
	wx.MustWriteString(userId, fileId, target, "0")
	result := hash.ToHex(digester)
	// Second pass
	digester.Reset()
	wx.MustWriteString(userKey, result, "000000")
	return hash.ToHexUpper(digester)
}
