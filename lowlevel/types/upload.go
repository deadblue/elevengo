package types

import (
	"time"

	"github.com/deadblue/elevengo/internal/util"
)

type UploadInfoResult struct {
	UserId  int
	UserKey string
}

type UploadInitParams struct {
	// Timestamp in seconds
	Timestamp int64
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
	// Update token
	Token string
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

type UploadTokenResult struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Expiration      time.Time
}

type UploadSampleInitResult struct {
	Host        string
	Object      string
	Callback    string
	AccessKeyId string
	Policy      string
	Signature   string
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
