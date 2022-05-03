package elevengo

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/oss"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// UploadTicket contains all required information to upload a file.
type UploadTicket struct {
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
	params *webapi.UploadOssParams) (exist bool, err error) {
	if !a.ut.Available() {
		if err = a.uploadInitToken(); err != nil {
			return
		}
	}
	// Prepare request
	targetId := fmt.Sprintf("U_1_%s", dirId)
	qs := web.Params{}.
		With("appid", a.ut.AppId).
		With("appversion", a.ut.AppVer).
		WithInt("isp", a.ut.IspType).
		With("sig", a.uploadCalculateSignature(targetId, quickId)).
		With("format", "json").
		WithNow("t")
	form := web.Params{}.
		With("app_ver", a.ut.AppVer).
		With("preid", preId).
		With("quickid", quickId).
		With("target", targetId).
		With("fileid", quickId).
		With("filename", name).
		WithInt64("filesize", size).
		WithInt("userid", a.ut.UserId)
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
	wx.MustWriteString(strconv.Itoa(a.uid), fileId, fileId, targetId, "0")
	h := hex.EncodeToString(digester.Sum(nil))
	// Second pass
	digester.Reset()
	wx.MustWriteString(a.ut.UserKey, h, "000000")
	return strings.ToUpper(hex.EncodeToString(digester.Sum(nil)))
}

// UploadCreateTicket creates a ticket with all required information to upload a
// file. Caller can use third-party tools/libraries to process it.
func (a *Agent) UploadCreateTicket(dirId, name string, r io.Reader, ticket *UploadTicket) (exist bool, err error) {
	// Digest content
	dr := &hash.DigestResult{}
	if err = hash.Digest(r, dr); err != nil {
		return
	}
	// Initialize uploading
	params := &webapi.UploadOssParams{}
	if exist, err = a.uploadInit(dirId, name, dr.Size, dr.PreId, dr.QuickId, params); exist || err != nil {
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
func (a *Agent) UploadParseResult(content []byte, file *File) (err error) {
	resp := &webapi.BasicResponse{}
	if err = json.Unmarshal(content, resp); err == nil {
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
	file.Size = data.FileSize
	file.Sha1 = data.Sha1
	file.PickCode = data.PickCode
	return
}

func (a *Agent) UploadSimply(dirId, name string, size int64, r io.Reader) (err error) {
	if size == 0 {
		// Try to inspect input size from r
		switch r.(type) {
		case *bytes.Buffer:
			size = int64(r.(*bytes.Buffer).Len())
		case *bytes.Reader:
			size = r.(*bytes.Reader).Size()
		case *strings.Reader:
			size = int64(r.(*strings.Reader).Len())
		case *os.File:
			if i, e := r.(*os.File).Stat(); e == nil {
				size = i.Size()
			}
		}
	}
	if size == 0 {
		return errors.New("upload size is zero")
	}
	form := web.Params{}.
		WithInt("userid", a.ut.UserId).
		With("filename", name).
		WithInt64("filesize", size).
		With("target", fmt.Sprintf("U_1_%s", dirId))
	resp := &webapi.UploadSimpleInitResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadSimpleInit, nil, form, resp); err != nil {
		return
	}
	// TODO: Unfinished

	return
}
