package webapi

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
)

const (
	UploadMaxSize = 5 * 1024 * 1024 * 1024

	UploadSimplyMaxSize = 200 * 1024 * 1024

	UploadPreSize = 128 * 1024

	uploadTokenPrefix = "Qclm8MGWUv59TnrR0XPg"
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
	ErrorCode int    `json:"statuscode"`
	ErrorMsg  string `json:"statusmsg"`

	Status   BoolInt `json:"status"`
	PickCode string  `json:"pickcode"`

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
	if r.ErrorCode == 0 {
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

func UploadCalculateToken(
	userId, fileId, preId string, 
	fileSize, timestamp int64,
) string {
	userHash := hash.Md5Hex(userId)
	digester := md5.New()
	wx := util.UpgradeWriter(digester)
	wx.MustWriteString(
		uploadTokenPrefix,
		fileId, 
		strconv.FormatInt(fileSize, 10), 
		preId,
		userId, 
		strconv.FormatInt(timestamp, 10), 
		userHash, 
		AppVersion)
	return hex.EncodeToString(digester.Sum(nil))
}

func UploadCalculateSignature(userId, userKey, fileId, targetId string) string {
	digester := sha1.New()
	wx := util.UpgradeWriter(digester)
	// First pass
	wx.MustWriteString(userId, fileId, targetId, "0")
	h := hex.EncodeToString(digester.Sum(nil))
	// Second pass
	digester.Reset()
	wx.MustWriteString(userKey, h, "000000")
	return strings.ToUpper(hex.EncodeToString(digester.Sum(nil)))
}
