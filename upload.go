package elevengo

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/oss"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/types"
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

type _UploadOssParams struct {
	// File size
	Size int64
	// Base64-encoded MD5 hash value
	MD5 string
	// Bucket name on OSS
	Bucket string
	// Object name on OSS
	Object string
	// Callback parameters
	Callback    string
	CallbackVar string
}

func (a *Agent) uploadInit(
	dirId, name string,
	rs io.ReadSeeker, maxSize int64,
	op *_UploadOssParams,
) (exists bool, err error) {
	// Digest incoming stream
	dr := &hash.DigestResult{}
	if err = hash.Digest(rs, dr); err != nil {
		return
	}
	if maxSize > 0 && dr.Size > maxSize {
		err = errors.ErrUploadTooLarge
		return
	}
	var signKey, signValue string
	// Call API
	for {
		spec := (&api.UploadInitSpec{}).Init(
			dirId, dr.SHA1, name, dr.Size, signKey, signValue, &a.common,
		)
		if err = a.llc.CallApi(spec, context.Background()); err != nil {
			break
		}
		if spec.Result.SignKey != "" {
			// Update parameters
			signKey = spec.Result.SignKey
			signValue, _ = hash.DigestRange(rs, spec.Result.SignCheck)
		} else {
			if spec.Result.Exists {
				exists = true
			} else if op != nil {
				op.Size, op.MD5 = dr.Size, dr.MD5
				op.Bucket = spec.Result.Oss.Bucket
				op.Object = spec.Result.Oss.Object
				op.Callback = spec.Result.Oss.Callback
				op.CallbackVar = spec.Result.Oss.CallbackVar
			}
			break
		}
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
	dirId, name string, r io.ReadSeeker, ticket *UploadTicket,
) (err error) {
	// Initialize uploading
	op := &_UploadOssParams{}
	if ticket.Exist, err = a.uploadInit(
		dirId, name, r, api.UploadMaxSize, op,
	); err != nil || ticket.Exist {
		return
	}
	// Get OSS token
	spec := (&api.UploadTokenSpec{}).Init()
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}
	// Fill UploadTicket
	ticket.Expiration = spec.Result.Expiration
	ticket.Verb = http.MethodPut
	ticket.Url = oss.GetPutObjectUrl(op.Bucket, op.Object)
	if ticket.Header == nil {
		ticket.Header = make(map[string]string)
	}
	ticket.header(oss.HeaderDate, oss.Date()).
		header(oss.HeaderContentLength, strconv.FormatInt(op.Size, 10)).
		header(oss.HeaderContentType, util.DetermineMimeType(name)).
		header(oss.HeaderContentMd5, op.MD5).
		header(oss.HeaderOssCallback, util.Base64Encode(op.Callback)).
		header(oss.HeaderOssCallbackVar, util.Base64Encode(op.CallbackVar)).
		header(oss.HeaderOssSecurityToken, spec.Result.SecurityToken)

	authorization := oss.CalculateAuthorization(&oss.RequestMetadata{
		Verb:   ticket.Verb,
		Header: ticket.Header,
		Bucket: op.Bucket,
		Object: op.Object,
	}, spec.Result.AccessKeyId, spec.Result.AccessKeySecret)
	ticket.header(oss.HeaderAuthorization, authorization)
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
	dirId, name string, r io.ReadSeeker, ticket *UploadOssTicket,
) (err error) {
	// Get OSS parameters
	op := &_UploadOssParams{}
	if ticket.Exist, err = a.uploadInit(
		dirId, name, r, -1, op,
	); err != nil || ticket.Exist {
		return
	}
	// Get OSS token
	spec := (&api.UploadTokenSpec{}).Init()
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}
	// Fill ticket
	ticket.Expiration = spec.Result.Expiration
	ticket.Client.Endpoint = oss.GetEndpointUrl()
	ticket.Client.AccessKeyId = spec.Result.AccessKeyId
	ticket.Client.AccessKeySecret = spec.Result.AccessKeySecret
	ticket.Client.SecurityToken = spec.Result.SecurityToken
	ticket.Bucket = op.Bucket
	ticket.Object = op.Object
	ticket.Callback = util.Base64Encode(op.Callback)
	ticket.CallbackVar = util.Base64Encode(op.CallbackVar)
	return
}

// UploadParseResult parses the raw upload response, and fills it to file.
func (a *Agent) UploadParseResult(r io.Reader, file *File) (err error) {
	jd, resp := json.NewDecoder(r), &protocol.StandardResp{}
	if err = jd.Decode(resp); err == nil {
		err = resp.Err()
	}
	if err != nil || file == nil {
		return
	}
	result := &types.UploadSampleResult{}
	if err = resp.Extract(result); err != nil {
		return
	}
	// Note: Not all fields of file are filled.
	file.IsDirectory = false
	file.FileId = result.FileId
	file.Name = result.FileName
	file.Size = result.FileSize.Int64()
	file.Sha1 = result.FileSha1
	file.PickCode = result.PickCode
	return
}

// UploadSample directly uploads small file/data (smaller than 200MB) to cloud.
func (a *Agent) UploadSample(dirId, name string, size int64, r io.Reader) (fileId string, err error) {
	if size == 0 {
		size = util.GuessSize(r)
	}
	// Check upload size
	if size <= 0 {
		// What the fuck?
		return "", errors.ErrUploadNothing
	} else if size > api.UploadMaxSizeSample {
		return "", errors.ErrUploadTooLarge
	}
	// Call API.
	initSpec := (&api.UploadSampleInitSpec{}).Init(dirId, name, size, &a.common)
	if err = a.llc.CallApi(initSpec, context.Background()); err != nil {
		return
	}
	// Upload file
	upSpec := (&api.UploadSampleSpec{}).Init(dirId, name, r, &initSpec.Result)
	if err = a.llc.CallApi(upSpec, context.Background()); err == nil {
		fileId = upSpec.Result.FileId
	}
	return
}
