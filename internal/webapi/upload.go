package webapi

import (
	"crypto/md5"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
)

const (
	UploadMaxSize = 5 * 1024 * 1024 * 1024
	UploadSimplyMaxSize = 200 * 1024 * 1024

	UploadStatusNormal = 1
	UploadStatusExist = 2
	UploadStatusDoubleCheck = 7

	uploadTokenSalt = "Qclm8MGWUv59TnrR0XPg"
)

type UploadToken struct {
	AppId   string
	AppVer  string
	IspType int
	UserId  int
	UserKey string
}

func (t *UploadToken) Available() bool {
	return t.UserKey != ""
}

type UploadResultData struct {
	AreaId     IntString   `json:"aid"`
	CategoryId IntString   `json:"cid"`
	FileId     string      `json:"file_id"`
	FileName   string      `json:"file_name"`
	FileSize   StringInt64 `json:"file_size"`
	PickCode   string      `json:"pick_code"`
	Sha1       string      `json:"sha1"`
}

type UploadInfoResponse struct {
	BasicResponse
	AppId       IntString `json:"app_id"`
	AppVersion  IntString `json:"app_version"`
	UploadLimit int64     `json:"size_limit"`
	IspType     int       `json:"isp_type"`
	UserId      int       `json:"user_id"`
	UserKey     string    `json:"userkey"`
}

type UploadInitResponse struct {
	Request   string `json:"request"`
	Version   string `json:"version"`
	ErrorCode int    `json:"statuscode"`
	ErrorMsg  string `json:"statusmsg"`

	Status   BoolInt `json:"status"`
	PickCode string  `json:"pickcode"`

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

func (r *UploadInitResponse) Err() error {
	// Ignore 701 error
	if r.ErrorCode == 0 || r.ErrorCode == 701 {
		return nil
	}
	return errors.New(r.ErrorMsg)
}

type UploadOssParams struct {
	Bucket      string
	Object      string
	Callback    string
	CallbackVar string
}

type UploadOssTokenResponse struct {
	StatusCode      string `json:"StatusCode"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"Expiration"`
}

func (r *UploadOssTokenResponse) Err() error {
	if r.StatusCode == "200" {
		return nil
	}
	return ErrUnexpected
}

type UploadSimpleInitResponse struct {
	Host     string `json:"host"`
	Object   string `json:"object"`
	Callback string `json:"callback"`

	AccessKeyId string `json:"accessid"`
	Policy      string `json:"policy"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
}

func (r *UploadSimpleInitResponse) Err() error {
	return nil
}

// UploadHelper is a helper object for upload function.
type UploadHelper struct {
	userId   string
	userKey  string
	userHash string
}

func (h *UploadHelper) IsReady() bool {
	return h.userKey != ""
}

func (h *UploadHelper) UserId() string {
	return h.userId
}

func (h *UploadHelper) Init(userId int, userKey string) {
	h.userId = strconv.Itoa(userId)
	h.userKey = userKey
	// Calculate user hash only once
	h.userHash = hash.Md5Hex(h.userId)
}

func (h *UploadHelper) CalculateSignature(fileId, targetId string) string {
	digester := sha1.New()
	wx := util.UpgradeWriter(digester)
	// First pass
	wx.MustWriteString(h.userId, fileId, targetId, "0")
	result := hash.ToHex(digester)
	// Second pass
	digester.Reset()
	wx.MustWriteString(h.userKey, result, "000000")
	return hash.ToHexUpper(digester)
}

func (h *UploadHelper) CalculateToken(
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
		h.userId,
		strconv.FormatInt(timestamp, 10), 
		h.userHash,
		AppVersion,
	)
	return hash.ToHex(digester)
}

type UploadDigestResult struct {
	FileId   string
	FileSize int64
	MD5      string
}

func UploadDigest(r io.Reader, result *UploadDigestResult) (err error) {
	hs, hm := sha1.New(), md5.New()
	w := io.MultiWriter(hs, hm)
	// Write remain data.
	if result.FileSize, err = io.Copy(w, r); err != nil {
		return
	}
	result.FileId, result.MD5 = hash.ToHexUpper(hs), hash.ToBase64(hm)
	return nil
}

func UploadDigestRange(r io.ReadSeeker, rangeSpec string) (result string, err error) {
	var start, end int64
	if _, err = fmt.Sscanf(rangeSpec, "%d-%d", &start, &end); err != nil {
		return
	}
	h := sha1.New()
	r.Seek(start, io.SeekStart)
	if _, err = io.CopyN(h, r, end - start + 1); err == nil {
		result = hash.ToHexUpper(h)
	}
	return
}