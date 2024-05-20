package types

import (
	"time"

	"github.com/deadblue/elevengo/internal/util"
)

type UploadInfoResult struct {
	UserId  int
	UserKey string
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
	// Check parameters
	SignKey   string
	SignCheck string
	// Pickcode is available when rapid-uploaded
	PickCode string
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
