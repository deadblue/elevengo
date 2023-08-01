package api

import (
	"errors"

	"github.com/deadblue/elevengo/internal/api/base"
)

type _UploadInfoResp struct {
	base.BasicResp
	UserId  int
	UserKey string
}

type UploadInfoSpec struct {
	base.JsonApiSpec[_UploadInfoResp]
}

func (s *UploadInfoSpec) Init() *UploadInfoSpec {
	s.JsonApiSpec.Init("https://proapi.115.com/app/uploadinfo")
	return s
}

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
	return errors.New(r.ErrorMsg)
}

type UploadInitSpec struct {
	base.JsonApiSpec[_UploadInitResp]
}

func (s *UploadInitSpec) Init() *UploadInitSpec {
	s.JsonApiSpec.Init("https://uplb.115.com/4.0/initupload.php")
	s.JsonApiSpec.EnableCrypto()
	return s
}
