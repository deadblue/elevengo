package api

import (
	"crypto/md5"
	"crypto/sha1"
	"io"
	"strconv"
	"time"

	"github.com/deadblue/elevengo/internal/api/base"
	"github.com/deadblue/elevengo/internal/api/errors"
	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
)

const (
	uploadTokenSalt = "Qclm8MGWUv59TnrR0XPg"
)

type UploadInfoResult struct {
	UserId  int
	UserKey string
}

//lint:ignore U1000 This type is used in generic.
type _UploadInfoResp struct {
	base.BasicResp
	UserId  int    `json:"user_id"`
	UserKey string `json:"userkey"`
}

func (r *_UploadInfoResp) Extract(v any) error {
	if ptr, ok := v.(*UploadInfoResult); !ok {
		return errors.ErrUnsupportedResult
	} else {
		ptr.UserId = r.UserId
		ptr.UserKey = r.UserKey
	}
	return nil
}

type UploadInfoSpec struct {
	base.JsonApiSpec[UploadInfoResult, _UploadInfoResp]
}

func (s *UploadInfoSpec) Init() *UploadInfoSpec {
	s.JsonApiSpec.Init("https://proapi.115.com/app/uploadinfo")
	return s
}

type UploadHelper struct {
	AppVer string
	UserId string

	userHash string
	userKey  string
}

func (h *UploadHelper) SetUserParams(userId int, userKey string) {
	h.UserId = strconv.Itoa(userId)
	h.userKey = userKey
	h.userHash = hash.Md5Hex(h.UserId)
}

func (h *UploadHelper) CalcSign(fileId, target string) string {
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

func (h *UploadHelper) CalcToken(
	fileId string, fileSize int64,
	signKey, signValue string,
	timestamp int64,
) string {
	digester := md5.New()
	wx := util.UpgradeWriter(digester)
	wx.MustWriteString(
		uploadTokenSalt,
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

//lint:ignore U1000 This type is used in generic.
type _UploadInitResp struct {
	Request   string `json:"request"`
	Version   string `json:"version"`
	ErrorCode int    `json:"statuscode"`
	ErrorMsg  string `json:"statusmsg"`

	Status   int    `json:"status"`
	PickCode string `json:"pickcode"`

	// New fields in upload v4.0
	SignKey   string `json:"sign_key"`
	SignCheck string `json:"sign_check"`

	// OSS upload fields
	Bucket   string `json:"bucket"`
	Object   string `json:"object"`
	Callback struct {
		Callback    string `json:"callback"`
		CallbackVar string `json:"callback_var"`
	} `json:"callback"`

	// Useless fields
	FileId   int    `json:"fileid"`
	FileInfo string `json:"fileinfo"`
	Target   string `json:"target"`
}

func (r *_UploadInitResp) Err() error {
	// Ignore 701 error
	if r.ErrorCode == 0 || r.ErrorCode == 701 {
		return nil
	}
	return errors.ErrUnexpected
}

type UploadInitSpec struct {
	base.JsonApiSpec[any, _UploadInitResp]
}

func (s *UploadInitSpec) Init(uh *UploadHelper, r io.ReadSeeker) *UploadInitSpec {
	s.JsonApiSpec.Init("https://uplb.115.com/4.0/initupload.php")
	s.EnableCrypto()
	// Prepare parameters
	now := time.Now().Unix()
	s.FormSetAll(map[string]string{
		"appid":      "0",
		"appversion": uh.AppVer,
		"userid":     uh.UserId,
	})
	s.FormSetInt64("t", now)
	return s
}
