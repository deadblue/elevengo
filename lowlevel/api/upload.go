package api

import (
	"io"

	"github.com/deadblue/elevengo/internal/multipart"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/types"
)

const (
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

type UploadInitSpec struct {
	_JsonApiSpec[types.UploadInitResult, protocol.UploadInitResp]
}

func (s *UploadInitSpec) Init(params *types.UploadInitParams, userId string, appVer string) *UploadInitSpec {
	s._JsonApiSpec.Init("https://uplb.115.com/4.0/initupload.php")
	s.crypto = true
	// Prepare parameters
	s.form.Set("appid", "0").
		Set("appversion", appVer).
		Set("userid", userId).
		Set("filename", params.FileName).
		SetInt64("filesize", params.FileSize).
		Set("fileid", params.FileId).
		Set("target", params.Target).
		Set("sig", params.Signature).
		SetInt64("t", params.Timestamp).
		Set("token", params.Token)
	if params.SignKey != "" && params.SignValue != "" {
		s.form.Set("sign_key", params.SignKey).
			Set("sign_val", params.SignValue)
	}
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

func (s *UploadSampleSpec) Init(
	name string, target string, r io.Reader,
	initResult *types.UploadSampleInitResult,
) *UploadSampleSpec {
	s._StandardApiSpec.Init(initResult.Host)
	//Prepart payload
	s.payload = multipart.Builder().
		AddValue("success_action_status", "200").
		AddValue("name", name).
		AddValue("target", target).
		AddValue("key", initResult.Object).
		AddValue("policy", initResult.Policy).
		AddValue("OSSAccessKeyId", initResult.AccessKeyId).
		AddValue("callback", initResult.Callback).
		AddValue("signature", initResult.Signature).
		AddFile("file", name, r).
		Build()
	return s
}

func (s *UploadSampleSpec) Payload() client.Payload {
	return s.payload
}
