package elevengo

import (
	"encoding/json"
	"fmt"
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"github.com/deadblue/gostream/multipart"
	"github.com/deadblue/gostream/quietly"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	apiUploadInit = "https://uplb.115.com/3.0/sampleinitupload.php"
)

/*
UploadInfo contains all required information to create an upload ticket.

To upload a regular file, caller can use os.FileInfo as UploadInfo, see
"Agent.UploadCreateTicket" doc for detail.

To upload a memory data as file, caller should implement it himself.
*/
type UploadInfo interface {
	// Name of the upload file.
	Name() string
	// Size in bytes of the upload file.
	Size() int64
}

// UploadTicket contains all required information to upload a file.
type UploadTicket struct {
	// Remote URL which will receive the file content.
	Endpoint string
	// Field name for the upload file.
	FileField string
	// Other parameters that should be sent with the file.
	Values map[string]string
}

/*
UploadCreateTicket creates a ticket which contains all required information
to upload a file. Caller can use third-party tools/libraries to upload file.
*/
func (a *Agent) UploadCreateTicket(parentId string, info UploadInfo) (ticket UploadTicket, err error) {
	// Request upload token
	form := core.NewForm().
		WithInt("userid", a.ui.Id).
		WithString("filename", info.Name()).
		WithInt64("filesize", info.Size()).
		WithString("target", fmt.Sprintf("U_1_%s", parentId))
	result := &types.UploadInitResult{}
	if err = a.hc.JsonApi(apiUploadInit, nil, form, result); err != nil {
		return
	}
	// Create upload ticket
	ticket = UploadTicket{
		Endpoint:  result.Host,
		FileField: "file",
		Values: map[string]string{
			"OSSAccessKeyId": result.AccessKeyId,
			"key":            result.ObjectKey,
			"policy":         result.Policy,
			"callback":       result.Callback,
			"signature":      result.Signature,
			"name":           info.Name(),
		},
	}
	return
}

/*
UploadParseResult parses the raw upload response body to file metadata. It
is useful when caller process upload ticket through thirdpary tools/libraries.
*/
func (a *Agent) UploadParseResult(content []byte) (file *File, err error) {
	result := &types.UploadResult{}
	if err = json.Unmarshal(content, result); err == nil {
		data := result.Data
		createTime := time.Unix(data.CreateTime, 0)
		file = &File{
			IsFile:      true,
			IsDirectory: false,
			FileId:      data.FileId,
			ParentId:    data.CategoryId,
			Name:        data.FileName,
			Size:        int64(data.FileSize),
			PickCode:    data.PickCode,
			Sha1:        data.Sha1,
			CreateTime:  createTime,
			UpdateTime:  createTime,
		}
	}
	return
}

/*
Upload uploads data as a file to cloud, and returns the file metadata on
success. If r implements io.Closer, it will be closed automatically.
*/
func (a *Agent) Upload(parentId string, info UploadInfo, r io.Reader) (file *File, err error) {
	// Register defer function only when r implements io.Closer.
	if rc, ok := r.(io.ReadCloser); ok {
		defer quietly.Close(rc)
	}

	ticket, err := a.UploadCreateTicket(parentId, info)
	if err != nil {
		return
	}
	// Create multipart form for uploading
	form := multipart.New()
	for name, value := range ticket.Values {
		form.AddValue(name, value)
	}
	form.AddFileData(ticket.FileField, info.Name(), info.Size(), r)
	// Make request
	req, err := multipart.NewRequest(ticket.Endpoint, form)
	if err != nil {
		return
	}
	// Send request through default client
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer quietly.Close(resp.Body)
	// Parse response body
	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		return a.UploadParseResult(body)
	}
}
