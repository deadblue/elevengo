package upload

import (
	"crypto/md5"
	"strconv"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
)

const (
	tokenSalt = "Qclm8MGWUv59TnrR0XPg"
)

func CalcToken(
	appVer, userId, userHash string,
	fileId string, fileSize int64,
	signKey, signValue string,
	timestamp int64,
) string {
	digester := md5.New()
	wx := util.UpgradeWriter(digester)
	wx.MustWriteString(
		tokenSalt,
		fileId,
		strconv.FormatInt(fileSize, 10),
		signKey,
		signValue,
		userId,
		strconv.FormatInt(timestamp, 10),
		userHash,
		appVer,
	)
	return hash.ToHex(digester)
}
