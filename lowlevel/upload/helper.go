package upload

import (
	"crypto/md5"
	"crypto/sha1"
	"strconv"
	"time"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/types"
)

const (
	tokenSalt = "Qclm8MGWUv59TnrR0XPg"
)

type Helper struct {
	AppVer string
	UserId string

	userHash string
	userKey  string
}

func (h *Helper) SetUserParams(userId int, userKey string) {
	h.UserId = strconv.Itoa(userId)
	h.userKey = userKey
	h.userHash = hash.Md5Hex(h.UserId)
}

func (h *Helper) calcSign(fileId, target string) string {
	digester := sha1.New()
	wx := util.UpgradeWriter(digester)
	// First pass
	wx.MustWriteString(h.UserId, fileId, target, "0")
	result := hash.ToHex(digester)
	// Second pass
	digester.Reset()
	wx.MustWriteString(h.userKey, result, "000000")
	return hash.ToHexUpper(digester)
}

func (h *Helper) calcToken(
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
		h.UserId,
		strconv.FormatInt(timestamp, 10),
		h.userHash,
		h.AppVer,
	)
	return hash.ToHex(digester)
}

func (h *Helper) Sign(params *types.UploadInitParams) *types.UploadInitParams {
	if params.Timestamp == 0 {
		params.Timestamp = time.Now().Unix()
	}
	if params.Signature == "" {
		params.Signature = h.calcSign(
			params.FileId, params.Target,
		)
	}
	params.Token = h.calcToken(
		params.FileId, params.FileSize,
		params.SignKey, params.SignValue,
		params.Timestamp,
	)
	return params
}
