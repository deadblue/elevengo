package protocol

import (
	"time"

	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/types"
)

//lint:ignore U1000 This type is used in generic.
type UploadInfoResp struct {
	BasicResp

	UserId  int    `json:"user_id"`
	UserKey string `json:"userkey"`
}

func (r *UploadInfoResp) Extract(v any) error {
	if ptr, ok := v.(*types.UploadInfoResult); !ok {
		return errors.ErrUnsupportedResult
	} else {
		ptr.UserId = r.UserId
		ptr.UserKey = r.UserKey
	}
	return nil
}

//lint:ignore U1000 This type is used in generic.
type UploadInitResp struct {
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

func (r *UploadInitResp) Err() error {
	// Ignore 701 error
	if r.ErrorCode == 0 || r.ErrorCode == 701 {
		return nil
	}
	return errors.ErrUnexpected
}

func (r *UploadInitResp) Extract(v any) (err error) {
	if ptr, ok := v.(*types.UploadInitResult); !ok {
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

//lint:ignore U1000 This type is used in generic.
type UploadTokenResp struct {
	StatusCode      string `json:"StatusCode"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"Expiration"`
}

func (r *UploadTokenResp) Err() error {
	if r.StatusCode == "200" {
		return nil
	}
	return errors.ErrUnexpected
}

func (r *UploadTokenResp) Extract(v any) error {
	if ptr, ok := v.(*types.UploadTokenResult); !ok {
		return errors.ErrUnsupportedResult
	} else {
		ptr.AccessKeyId = r.AccessKeyId
		ptr.AccessKeySecret = r.AccessKeySecret
		ptr.SecurityToken = r.SecurityToken
		ptr.Expiration, _ = time.Parse(time.RFC3339, r.Expiration)
	}
	return nil
}

//lint:ignore U1000 This type is used in generic.
type UploadSampleInitResp struct {
	Host        string `json:"host"`
	Object      string `json:"object"`
	Callback    string `json:"callback"`
	AccessKeyId string `json:"accessid"`
	Policy      string `json:"policy"`
	Signature   string `json:"signature"`
	Expire      int64  `json:"expire"`
}

func (r *UploadSampleInitResp) Err() error {
	return nil
}

func (r *UploadSampleInitResp) Extract(v any) error {
	if ptr, ok := v.(*types.UploadSampleInitResult); !ok {
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
