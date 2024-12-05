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

func (r *UploadInfoResp) Extract(v *types.UploadInfoResult) error {
	v.UserId = r.UserId
	v.UserKey = r.UserKey
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

func (r *UploadInitResp) Extract(v *types.UploadInitResult) (err error) {
	switch r.Status {
	case 1:
		v.Oss.Bucket = r.Bucket
		v.Oss.Object = r.Object
		v.Oss.Callback = r.Callback.Callback
		v.Oss.CallbackVar = r.Callback.CallbackVar
	case 2:
		v.Exists = true
		v.PickCode = r.PickCode
	case 7:
		v.SignKey = r.SignKey
		v.SignCheck = r.SignCheck
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

func (r *UploadTokenResp) Extract(v *types.UploadTokenResult) error {
	v.AccessKeyId = r.AccessKeyId
	v.AccessKeySecret = r.AccessKeySecret
	v.SecurityToken = r.SecurityToken
	v.Expiration, _ = time.Parse(time.RFC3339, r.Expiration)
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

func (r *UploadSampleInitResp) Extract(v *types.UploadSampleInitResult) error {
	v.Host = r.Host
	v.Object = r.Object
	v.Callback = r.Callback
	v.AccessKeyId = r.AccessKeyId
	v.Policy = r.Policy
	v.Signature = r.Signature
	return nil
}
