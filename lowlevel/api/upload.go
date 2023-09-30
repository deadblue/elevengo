package api

import (
	"crypto/md5"
	"crypto/sha1"
	"strconv"
	"time"

	"github.com/deadblue/elevengo/internal/apibase"
	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/protocol"
)

const (
	uploadTokenSalt = "Qclm8MGWUv59TnrR0XPg"

	UploadMaxSize       = 5 * 1024 * 1024 * 1024
	UploadMaxSizeSample = 200 * 1024 * 1024
)

type UploadInfoResult struct {
	UserId  int
	UserKey string
}

//lint:ignore U1000 This type is used in generic.
type _UploadInfoResp struct {
	apibase.BasicResp
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
	apibase.JsonApiSpec[UploadInfoResult, _UploadInfoResp]
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

type UploadInitResult struct {
	Exists bool
	// Upload parameters
	Oss struct {
		Bucket      string
		Object      string
		Callback    string
		CallbackVar string
	}
	SignKey   string
	SignCheck string
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

func (r *_UploadInitResp) Extract(v any) (err error) {
	if ptr, ok := v.(*UploadInitResult); !ok {
		err = errors.ErrUnsupportedResult
	} else {
		switch r.Status {
		case 1:
			ptr.Oss.Bucket = r.Bucket
			ptr.Oss.Object = r.Object
			ptr.Oss.Callback = r.Callback.Callback
			ptr.Oss.CallbackVar = r.Callback.CallbackVar
		case 2:
			ptr.Exists = true
		case 7:
			ptr.SignKey = r.SignKey
			ptr.SignCheck = r.SignCheck
		}
	}
	return
}

type UploadInitSpec struct {
	apibase.JsonApiSpec[UploadInitResult, _UploadInitResp]
}

func (s *UploadInitSpec) Init(params *UploadInitParams, helper *UploadHelper) *UploadInitSpec {
	s.JsonApiSpec.Init("https://uplb.115.com/4.0/initupload.php")
	s.EnableCrypto()
	// Prepare parameters
	now := time.Now().Unix()
	s.FormSet("appid", "0")
	s.FormSet("appversion", helper.AppVer)
	s.FormSet("userid", helper.UserId)
	s.FormSet("filename", params.FileName)
	s.FormSetInt64("filesize", params.FileSize)
	s.FormSet("fileid", params.FileId)
	s.FormSet("target", params.Target)
	s.FormSet("sig", params.Signature)
	s.FormSetInt64("t", now)
	if params.SignKey != "" && params.SignValue != "" {
		s.FormSet("sign_key", params.SignKey)
		s.FormSet("sign_val", params.SignValue)
	}
	s.FormSet("token", helper.CalcToken(
		params.FileId, params.FileSize, params.SignKey, params.SignValue, now,
	))
	return s
}

type UploadTokenResult struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      time.Time
}

//lint:ignore U1000 This type is used in generic.
type _UploadTokenResp struct {
	StatusCode      string `json:"StatusCode"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"Expiration"`
}

func (r *_UploadTokenResp) Err() error {
	if r.StatusCode == "200" {
		return nil
	}
	return errors.ErrUnexpected
}

func (r *_UploadTokenResp) Extract(v any) error {
	if ptr, ok := v.(*UploadTokenResult); !ok {
		return errors.ErrUnsupportedResult
	} else {
		ptr.AccessKeyId = r.AccessKeyId
		ptr.AccessKeySecret = r.AccessKeySecret
		ptr.SecurityToken = r.SecurityToken
		ptr.Expiration, _ = time.Parse(time.RFC3339, r.Expiration)
	}
	return nil
}

type UploadTokenSpec struct {
	apibase.JsonApiSpec[UploadTokenResult, _UploadTokenResp]
}

func (s *UploadTokenSpec) Init() *UploadTokenSpec {
	s.JsonApiSpec.Init("https://uplb.115.com/3.0/gettoken.php")
	return s
}

type UploadSampleInitResult struct {
	Host        string
	Object      string
	Callback    string
	AccessKeyId string
	Policy      string
	Signature   string
}

//lint:ignore U1000 This type is used in generic.
type _UploadSampleInitResp struct {
	Host        string `json:"host"`
	Object      string `json:"object"`
	Callback    string `json:"callback"`
	AccessKeyId string `json:"accessid"`
	Policy      string `json:"policy"`
	Signature   string `json:"signature"`
	Expire      int64  `json:"expire"`
}

func (r *_UploadSampleInitResp) Err() error {
	return nil
}

func (r *_UploadSampleInitResp) Extract(v any) error {
	if ptr, ok := v.(*UploadSampleInitResult); !ok {
		return errors.ErrUnsupportedResult
	} else {
		ptr.Host = r.Host
		ptr.Object = r.Object
		ptr.Callback = r.Callback
		ptr.AccessKeyId = r.AccessKeyId
		ptr.Policy = r.Policy
		ptr.Signature = r.Signature
	}
	return nil
}

type UploadSampleInitSpec struct {
	apibase.JsonApiSpec[UploadSampleInitResult, _UploadSampleInitResp]
}

func (s *UploadSampleInitSpec) Init(userId string, fileName string, fileSize int64, target string) *UploadSampleInitSpec {
	s.JsonApiSpec.Init("https://uplb.115.com/3.0/sampleinitupload.php")
	s.FormSet("userid", userId)
	s.FormSet("filename", fileName)
	s.FormSetInt64("filesize", fileSize)
	s.FormSet("target", target)
	return s
}

type UploadSampleResult struct {
	AreaId     util.IntNumber `json:"aid"`
	CategoryId string         `json:"cid"`
	FileId     string         `json:"file_id"`
	FileName   string         `json:"file_name"`
	FileSize   util.IntNumber `json:"file_size"`
	FileSha1   string         `json:"sha1"`
	PickCode   string         `json:"pick_code"`
	CreateTime util.IntNumber `json:"file_ptime"`
}

type UploadSampleSpec struct {
	apibase.JsonApiSpec[UploadSampleResult, apibase.StandardResp]
	payload protocol.Payload
}

func (s *UploadSampleSpec) Init(url string, payload protocol.Payload) *UploadSampleSpec {
	s.JsonApiSpec.Init(url)
	s.payload = payload
	return s
}

func (s *UploadSampleSpec) Payload() protocol.Payload {
	return s.payload
}
