package api

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/deadblue/elevengo/internal/multipart"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/upload"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/types"
)

const (
	UploadMaxSizeSample = 200 * 1024 * 1024
	UploadMaxSize       = 5 * 1024 * 1024 * 1024
	UploadMaxSizeOss    = 115 * 1024 * 1024 * 1024
)

func getTarget(dirId string) string {
	return fmt.Sprintf("U_1_%s", dirId)
}

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

func (s *UploadInitSpec) Init(
	dirId string, fileSha1 string, fileName string, fileSize int64,
	signKey string, signValue string,
	common *types.CommonParams,
) *UploadInitSpec {
	s._JsonApiSpec.Init("https://uplb.115.com/4.0/initupload.php")
	s.crypto = true
	// Make sure fileSha1 and signValue are in upper-case.
	fileSha1 = strings.ToUpper(fileSha1)
	signValue = strings.ToUpper(signValue)
	// Prepare parameters
	target := getTarget(dirId)
	timestamp := time.Now().UnixMilli()
	signature := upload.CalcSignature(common.UserId, common.UserKey, fileSha1, target)
	token := upload.CalcToken(
		common.AppVer, common.UserId, common.UserHash,
		fileSha1, fileSize, signKey, signValue, timestamp,
	)
	s.form.Set("appid", "0").
		Set("appversion", common.AppVer).
		Set("userid", common.UserId).
		Set("filename", fileName).
		SetInt64("filesize", fileSize).
		Set("fileid", fileSha1).
		Set("target", target).
		Set("sig", signature).
		SetInt64("t", timestamp).
		Set("token", token)
	if signKey != "" && signValue != "" {

		s.form.Set("sign_key", signKey).
			Set("sign_val", signValue)
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

func (s *UploadSampleInitSpec) Init(
	dirId string, fileName string, fileSize int64,
	common *types.CommonParams,
) *UploadSampleInitSpec {
	s._JsonApiSpec.Init("https://uplb.115.com/3.0/sampleinitupload.php")
	s.form.Set("userid", common.UserId).
		Set("filename", fileName).
		SetInt64("filesize", fileSize).
		Set("target", getTarget(dirId))
	return s
}

type UploadSampleSpec struct {
	_StandardApiSpec[types.UploadSampleResult]
	payload client.Payload
}

func (s *UploadSampleSpec) Init(
	dirId, fileName string, r io.Reader,
	initResult *types.UploadSampleInitResult,
) *UploadSampleSpec {
	s._StandardApiSpec.Init(initResult.Host)
	//Prepart payload
	s.payload = multipart.Builder().
		AddValue("success_action_status", "200").
		AddValue("name", fileName).
		AddValue("target", getTarget(dirId)).
		AddValue("key", initResult.Object).
		AddValue("policy", initResult.Policy).
		AddValue("OSSAccessKeyId", initResult.AccessKeyId).
		AddValue("callback", initResult.Callback).
		AddValue("signature", initResult.Signature).
		AddFile("file", fileName, r).
		Build()
	return s
}

func (s *UploadSampleSpec) Payload() client.Payload {
	return s.payload
}
