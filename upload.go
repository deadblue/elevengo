package elevengo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/deadblue/elevengo/internal/multipart"
	"github.com/deadblue/elevengo/internal/oss"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
)

// UploadTicket contains all required information to upload a file.
type UploadTicket struct {
	// Expiration time
	Expiration time.Time
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

func (a *Agent) uploadInitHelper() (err error) {
	if a.uh.IsReady() {
		return
	}
	resp := &webapi.UploadInfoResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadInfo, nil, nil, resp); err == nil {
		a.uh.Init(resp.UserId, resp.UserKey)
	}
	return
}

func (a *Agent) uploadInit(
	dirId, name string, 
	r io.ReadSeeker,
	dr *webapi.UploadDigestResult,
	op *webapi.UploadOssParams,
) (exist bool, err error) {
	if err = webapi.UploadDigest(r, dr); err != nil {
		return
	}
	if err = a.uploadInitHelper(); err != nil {
		return
	}
	// Prepare request parameters
	signKey, signVal := "", ""
	targetId := fmt.Sprintf("U_1_%s", dirId)
	form := web.Params{}.
		With("appid", "0").
		With("appversion", webapi.AppVersion).
		With("userid", a.uh.UserId()).
		With("filename", name).
		WithInt64("filesize", dr.FileSize).
		With("fileid", dr.FileId).
		With("target", targetId).
		With("sig", a.uh.CalculateSignature(dr.FileId, targetId))
	// Send request
	resp := &webapi.UploadInitResponse{}
	for retry := true; retry; {
		now := time.Now().Unix()
		form.WithInt64("t", now).
			With("token", a.uh.CalculateToken(dr.FileId, dr.FileSize, signKey, signVal, now))
		if signKey != "" && signVal != "" {
			form.With("sign_key", signKey).
				With("sign_val", signVal)
		}
		if err = a.wc.CallSecretJsonApi(
			webapi.ApiUploadInit, nil, form.ToForm(), resp, now,
		); err != nil {
			return
		}
		if resp.Status == webapi.UploadStatusDoubleCheck {
			// Update signKey & signVal
			signKey = resp.SignKey
			signVal, err = webapi.UploadDigestRange(r, resp.SignCheck)
			if err != nil {
				return
			}
		} else {
			retry = false
		}
	}
	// Parse response
	exist = resp.Status == webapi.UploadStatusExist
	if resp.Status == webapi.UploadStatusNormal && op != nil {
		op.Bucket = resp.Bucket
		op.Object = resp.Object
		op.Callback = resp.Callback.Callback
		op.CallbackVar = resp.Callback.CallbackVar
	}
	return
}

// UploadCreateTicket creates a ticket which contains all required parameters
// to upload file/data to cloud, the ticket should be used in 1 hour.
//
// To create ticket, r will be fully read to calculate SHA-1 and MD5 hash value. 
// If you want to re-use r, try to seek it to beginning.
// 
// To upload a file larger than 5G bytes, use `UploadCreateOssTicket`.
func (a *Agent) UploadCreateTicket(
	dirId, name string, r io.ReadSeeker, 
	ticket *UploadTicket,
) (err error) {
	// Initialize uploading
	dr, op := &webapi.UploadDigestResult{}, &webapi.UploadOssParams{}
	if ticket.Exist, err = a.uploadInit(
		dirId, name, r, dr, op,
	); ticket.Exist || err != nil {
		return
	}
	// Check file size
	if dr.FileSize > webapi.UploadMaxSize {
		return webapi.ErrUploadTooLarge
	}
	// Get OSS token
	resp := &webapi.UploadOssTokenResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadOssToken, nil, nil, resp); err != nil {
		return
	}
	// Fill UploadTicket
	ticket.Expiration, _ = time.Parse(time.RFC3339, resp.Expiration)
	ticket.Verb = http.MethodPut
	ticket.Url = oss.GetPutObjectUrl(op.Bucket, op.Object)
	if ticket.Header == nil {
		ticket.Header = make(map[string]string)
	}
	ticket.header(oss.HeaderDate, oss.Date()).
		header(oss.HeaderContentLength, strconv.FormatInt(dr.FileSize, 10)).
		header(oss.HeaderContentType, util.DetermineMimeType(name)).
		header(oss.HeaderContentMd5, dr.MD5).
		header(oss.HeaderOssCallback, util.Base64Encode(op.Callback)).
		header(oss.HeaderOssCallbackVar, util.Base64Encode(op.CallbackVar)).
		header(oss.HeaderOssSecurityToken, resp.SecurityToken)

	authorization := oss.CalculateAuthorization(&oss.RequestMetadata{
		Verb:   ticket.Verb,
		Header: ticket.Header,
		Bucket: op.Bucket,
		Object: op.Object,
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
		With("userid", a.uh.UserId()).
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

// UploadOssTicket contains all required paramters to upload a file through
// aliyun-oss-sdk(https://github.com/aliyun/aliyun-oss-go-sdk).
type UploadOssTicket struct {
	// Expiration time
	Expiration time.Time
	// Is file already exists
	Exist bool
	// Client parameters
	Client struct {
		Endpoint        string
		AccessKeyId     string
		AccessKeySecret string
		SecurityToken   string
	}
	// Bucket name
	Bucket string
	// Object key
	Object string
	// Callback option
	Callback string
	// CallbackVar option
	CallbackVar string
}

/*
UploadCreateOssTicket creates ticket to upload file through aliyun-oss-sdk. Use 
this method if you want to upload a file larger than 5G bytes.

To create ticket, r will be fully read to calculate SHA-1 and MD5 hash value. 
If you want to re-use r, try to seek it to beginning.

Example:

    import (
        "github.com/aliyun/aliyun-oss-go-sdk/oss"
        "github.com/deadblue/elevengo"
    )
    
	func main() {
		filePath := "/file/to/upload"

		var err error

		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Open file failed: %s", err)
		}
		defer file.Close()

		// Create 115 agent
		agent := elevengo.Default()
		if err = agent.CredentialImport(&elevengo.Credential{
			UID: "", CID: "", SEID: "",
		}); err != nil {
			log.Fatalf("Login failed: %s", err)
		}
		// Prepare OSS upload ticket
		ticket := &UploadOssTicket{}
		if err = agent.UploadCreateOssTicket(
			"dirId", 
			filepath.Base(file.Name()), 
			file, 
			ticket, 
		); err != nil {
			log.Fatalf("Create OSS ticket failed: %s", err)
		}
		if ticket.Exist {
			log.Printf("File has been fast-uploaded!")
			return
		}

		// Create OSS client
		oc, err := oss.New(
			ticket.Client.Endpoint, 
			ticket.Client.AccessKeyId,
			ticket.Client.AccessKeySecret,
			oss.SecurityToken(ticket.Client.SecurityToken)
		)
		if err != nil {
			log.Fatalf("Create OSS client failed: %s", err)
		}
		bucket, err := oc.Bucket(ticket.Bucket)
		if err != nil {
			log.Fatalf("Get OSS bucket failed: %s", err)
		}
		// Upload file in multipart.
		err = bucket.UploadFile(
			ticket.Object, 
			filePath, 
			100 * 1024 * 1024,	// 100 Megabytes per part
			oss.Callback(ticket.Callback),
			oss.CallbackVar(ticket.CallbackVar),
		)
		// Until now (2023-01-29), there is a bug in aliyun-oss-go-sdk:
		// When set Callback option, the response from CompleteMultipartUpload API 
		// is returned by callback host, which is not the standard XML. But SDK
		// always tries to parse it as CompleteMultipartUploadResult, and returns 
		// `io.EOF` error, just ignore it!
		if err != nil && err != io.EOF {
			log.Fatalf("Upload file failed: %s", err)
		} else {
			log.Print("Upload done!")
		}
	}
*/
func (a *Agent) UploadCreateOssTicket(
	dirId, name string, 
	r io.ReadSeeker, 
	ticket *UploadOssTicket,
) (err error) {
	// Initialize upload
	dr, op := &webapi.UploadDigestResult{}, &webapi.UploadOssParams{}
	if ticket.Exist, err = a.uploadInit(
		dirId, name, r, dr, op,
	); ticket.Exist || err != nil {
		return
	}
	// Get OSS token
	resp := &webapi.UploadOssTokenResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadOssToken, nil, nil, resp); err != nil {
		return
	}
	// Fill ticket
	ticket.Expiration, _ = time.Parse(time.RFC3339, resp.Expiration)
	ticket.Client.Endpoint = oss.GetEndpointUrl()
	ticket.Client.AccessKeyId = resp.AccessKeyId
	ticket.Client.AccessKeySecret = resp.AccessKeySecret
	ticket.Client.SecurityToken = resp.SecurityToken
	ticket.Bucket = op.Bucket
	ticket.Object = op.Object
	ticket.Callback = util.Base64Encode(op.Callback)
	ticket.CallbackVar = util.Base64Encode(op.CallbackVar)
	return
}
