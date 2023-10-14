package api

import (
	"crypto/md5"
	"crypto/sha1"
	"strconv"
	"time"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/types"
)

const (
	uploadTokenSalt = "Qclm8MGWUv59TnrR0XPg"

	UploadMaxSize       = 5 * 1024 * 1024 * 1024
	UploadMaxSizeSample = 200 * 1024 * 1024
)

type UploadInfoSpec struct {
	_JsonApiSpec[types.UploadInfoResult, protocol.UploadInfoResp]
}

func (s *UploadInfoSpec) Init() *UploadInfoSpec {
	s._JsonApiSpec.Init("https://proapi.115.com/app/uploadinfo")
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

type UploadInitParams struct {
	// File metadata
	FileId   string
	FileName string
	FileSize int64
	// Target directory
	Target string
	// Upload signature
	Signature string
	// Sign parameters for 2nd-pass
	SignKey   string
	SignValue string
}

type UploadInitSpec struct {
	_JsonApiSpec[types.UploadInitResult, protocol.UploadInitResp]
}

func (s *UploadInitSpec) Init(params *UploadInitParams, helper *UploadHelper) *UploadInitSpec {
	s._JsonApiSpec.Init("https://uplb.115.com/4.0/initupload.php")
	s.crypto = true
	// Prepare parameters
	now := time.Now().Unix()
	s.form.Set("appid", "0").
		Set("appversion", helper.AppVer).
		Set("userid", helper.UserId).
		Set("filename", params.FileName).
		SetInt64("filesize", params.FileSize).
		Set("fileid", params.FileId).
		Set("target", params.Target).
		Set("sig", params.Signature).
		SetInt64("t", now)
	if params.SignKey != "" && params.SignValue != "" {
		s.form.Set("sign_key", params.SignKey).
			Set("sign_val", params.SignValue)
	}
	s.form.Set("token", helper.CalcToken(
		params.FileId, params.FileSize, params.SignKey, params.SignValue, now,
	))
	return s
}

type UploadTokenSpec struct {
	_JsonApiSpec[types.UploadTokenResult, protocol.UploadTokenResp]
}

func (s *UploadTokenSpec) Init() *UploadTokenSpec {
	s._JsonApiSpec.Init("https://uplb.115.com/3.0/gettoken.php")
	return s
}

type UploadSampleInitSpec struct {
	_JsonApiSpec[types.UploadSampleInitResult, protocol.UploadSampleInitResp]
}

func (s *UploadSampleInitSpec) Init(userId string, fileName string, fileSize int64, target string) *UploadSampleInitSpec {
	s._JsonApiSpec.Init("https://uplb.115.com/3.0/sampleinitupload.php")
	s.form.Set("userid", userId).
		Set("filename", fileName).
		SetInt64("filesize", fileSize).
		Set("target", target)
	return s
}

type UploadSampleSpec struct {
	_StandardApiSpec[types.UploadSampleResult]
	payload client.Payload
}

func (s *UploadSampleSpec) Init(url string, payload client.Payload) *UploadSampleSpec {
	s._StandardApiSpec.Init(url)
	s.payload = payload
	return s
}

func (s *UploadSampleSpec) Payload() client.Payload {
	return s.payload
}
