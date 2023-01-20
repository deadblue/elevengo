package elevengo

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/multipart"
	"github.com/deadblue/elevengo/internal/oss"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// UploadTicket contains all required information to upload a file.
type UploadTicket struct {
	// Is file exists
	Exist bool
	// Request method
	Verb string
	// Remote URL which will receive the file content.
	Url string
	// Request header
	Header map[string]string
}

func (t *UploadTicket) header(name, value string) *UploadTicket {
	t.Header[name] = value
	return t
}

func (a *Agent) uploadInitToken() (err error) {
	resp := &webapi.UploadInfoResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadInfo, nil, nil, resp); err == nil {
		a.ut.AppId = string(resp.AppId)
		a.ut.AppVer = string(resp.AppVersion)
		a.ut.IspType = resp.IspType
		a.ut.UserId = resp.UserId
		a.ut.UserKey = resp.UserKey
	}
	return
}

func (a *Agent) uploadInit(
	dirId string, name string, size int64,
	preId string, quickId string,
	params *webapi.UploadOssParams,
) (exist bool, err error) {
	if !a.ut.Available() {
		if err = a.uploadInitToken(); err != nil {
			return
		}
	}
	// Prepare request
	now := time.Now().Unix()
	targetId := fmt.Sprintf("U_1_%s", dirId)
	qs := web.Params{}.
		With("appid", a.ut.AppId).
		With("appversion", webapi.AppVersion).
		WithInt("isp", a.ut.IspType).
		With("rt", "0").
		With("topupload", "0").
		With("token", a.uploadCalculateToken(quickId, size, preId, now)).
		With("sig", a.uploadCalculateSignature(targetId, quickId)).
		With("format", "json").
		WithInt64("t", now)
	form := web.Params{}.
		With("fileid", quickId).
		With("filename", name).
		WithInt64("filesize", size).
		With("preid", preId).
		With("target", targetId).
		WithInt("userid", a.ut.UserId).
		ToForm()
	// Send request
	resp := &webapi.UploadInitResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadInit, qs, form, resp); err != nil {
		return
	}
	// Parse response
	exist = resp.Status == 2
	if !exist && params != nil {
		params.Bucket = resp.Bucket
		params.Object = resp.Object
		params.Callback = resp.Callback.Callback
		params.CallbackVar = resp.Callback.CallbackVar
	}
	return
}

func (a *Agent) uploadCalculateSignature(targetId, fileId string) string {
	digester := sha1.New()
	wx := util.UpgradeWriter(digester)
	wx.MustWriteString(strconv.Itoa(a.ut.UserId), fileId, targetId, "0")
	h := hex.EncodeToString(digester.Sum(nil))
	// Second pass
	digester.Reset()
	wx.MustWriteString(a.ut.UserKey, h, "000000")
	return strings.ToUpper(hex.EncodeToString(digester.Sum(nil)))
}

func (a *Agent) uploadCalculateToken(
	fileId string, fileSize int64,
	preId string, timestamp int64,
) string {
	userId := strconv.Itoa(a.ut.UserId)
	userHash := hash.Md5Hex(userId)
	digester := md5.New()
	wx := util.UpgradeWriter(digester)
	wx.MustWriteString(webapi.UploadTokenPrefix,
		fileId, strconv.FormatInt(fileSize, 10), preId,
		userId, strconv.FormatInt(timestamp, 10), userHash,
		webapi.AppVersion)
	return hex.EncodeToString(digester.Sum(nil))
}

// UploadCreateTicket creates a ticket which contains all required information
// to upload file/data to cloud, the ticket should be used in 1 hour.
//
// To create ticket, r will be fully read to calculate SHA-1 hash and MD5
// hash. If you want to re-use r, try to seek it to beginning.
//
// Now, you can not upload file larger than 5GB, it will be supported later.
func (a *Agent) UploadCreateTicket(dirId, name string, r io.Reader, ticket *UploadTicket) (err error) {
	// Digest content
	dr := &hash.DigestResult{}
	if err = hash.Digest(r, dr); err != nil {
		return
	}
	if dr.Size > webapi.UploadMaxSize {
		return webapi.ErrUploadTooLarge
	}
	// Initialize uploading
	params := &webapi.UploadOssParams{}
	if ticket.Exist, err = a.uploadInit(
		dirId, name, dr.Size, dr.PreId, dr.QuickId, params); ticket.Exist || err != nil {
		return
	}

	// Get OSS token
	resp := &webapi.UploadOssTokenResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadOssToken, nil, nil, resp); err != nil {
		return
	}
	// Fill UploadTicket
	ticket.Verb = http.MethodPut
	ticket.Url = fmt.Sprintf("https://%s.%s/%s", params.Bucket, oss.Endpoint, params.Object)
	if ticket.Header == nil {
		ticket.Header = make(map[string]string)
	}
	ticket.header(oss.HeaderDate, oss.Date()).
		header(oss.HeaderContentLength, strconv.FormatInt(dr.Size, 10)).
		header(oss.HeaderContentType, util.DetermineMimeType(name)).
		header(oss.HeaderContentMd5, dr.MD5).
		header(oss.HeaderOssCallback, base64.StdEncoding.EncodeToString([]byte(params.Callback))).
		header(oss.HeaderOssCallbackVar, base64.StdEncoding.EncodeToString([]byte(params.CallbackVar))).
		header(oss.HeaderOssSecurityToken, resp.SecurityToken)

	authorization := oss.CalculateAuthorization(&oss.RequestMetadata{
		Verb:   ticket.Verb,
		Header: ticket.Header,
		Bucket: params.Bucket,
		Object: params.Object,
	}, resp.AccessKeyId, resp.AccessKeySecret)
	ticket.header(oss.HeaderAuthorization, authorization)
	return
}

// UploadParseResult parses the raw upload response, and fills it to file.
func (a *Agent) UploadParseResult(r io.Reader, file *File) (err error) {
	decoder, resp := json.NewDecoder(r), &webapi.BasicResponse{}
	if err = decoder.Decode(resp); err == nil {
		err = resp.Err()
	}
	if err != nil || file == nil {
		return
	}

	data := &webapi.UploadResultData{}
	if err = resp.Decode(data); err != nil {
		return
	}
	// Note: Not all fields of file are filled.
	file.IsDirectory = false
	file.FileId = data.FileId
	file.Name = data.FileName
	file.Size = int64(data.FileSize)
	file.Sha1 = data.Sha1
	file.PickCode = data.PickCode
	return
}

// UploadSimply directly uploads small file/data (smaller than 200MB) to cloud.
func (a *Agent) UploadSimply(dirId, name string, size int64, r io.Reader) (fileId string, err error) {
	if size == 0 {
		size = util.GuessSize(r)
	}
	// Check upload size
	if size <= 0 {
		// What the fuck?
		return "", errors.New("upload size is zero")
	} else if size > webapi.UploadSimplyMaxSize {
		return "", webapi.ErrUploadTooLarge
	}
	form := web.Params{}.
		WithInt("userid", a.ut.UserId).
		With("filename", name).
		WithInt64("filesize", size).
		With("target", fmt.Sprintf("U_1_%s", dirId)).
		ToForm()
	initResp := &webapi.UploadSimpleInitResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadSimpleInit, nil, form, initResp); err != nil {
		return
	}

	// Upload file
	mf := multipart.Builder().
		AddValue("success_action_status", "200").
		AddValue("name", name).
		AddValue("key", initResp.Object).
		AddValue("callback", initResp.Callback).
		AddValue("OSSAccessKeyId", initResp.AccessKeyId).
		AddValue("policy", initResp.Policy).
		AddValue("signature", initResp.Signature).
		AddFile("file", name, r).
		Build()
	uploadResp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(util.SecretUrl(initResp.Host), nil, mf, uploadResp); err != nil {
		return
	}
	// Parse response
	data := &webapi.UploadResultData{}
	if err = uploadResp.Decode(data); err == nil {
		fileId = data.FileId
	}
	return
}
