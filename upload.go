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
UploadCreateTicket creates a ticket which contails all necessary information to
upload a file. Caller can use thirdpary tools/libraries to perform the upload.

Example - Upload through curl:

	filename := "/path/to/file"
	// Get file info
	info, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	// Create ticket
	ticket, err := agent.UploadCreateTicket(parentId, info)
	if err != nil {
		log.Fatal(err)
	}
	// Create temp file to receive upload response
	tmpFile, err := ioutil.TempFile(os.TempDir(), "115-upload-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	// Construct curl command
	cmd := exec.Command("/path/to/curl", ticket.Endpoint, "-o", tmpFile.Name(), "-#")
	for name, value := range ticket.Values {
		cmd.Args = append(cmd.Args, "-F", fmt.Sprintf("%s=%s", name, value))
	}
	// NOTICE: File field should be at the end of the form.
	cmd.Args = append(cmd.Args, "-F", fmt.Sprintf("%s=@%s", ticket.FileField, filename))
	// Run the command
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
	// Parse upload response
	response, _ := ioutil.ReadAll(tmpFile)
	file, err := agent.UploadParseResult(response)
	if err != nil {
		log.Fatel(f)
	} else {
		log.Printf("Uploaded file: %#v", file)
	}
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

// UploadParseResult parses the raw upload response body. If caller performs upload ticket
// through thirdpary tools/libraries, he can call the method to parse the upload result.
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
Upload uploads data as a file to cloud, and returns the file metadata on successful.
If r implements io.Closer, it will be closed by method.

Example:

	// To upload a regular file
	file, err := os.Open("/path/to/file")
	if err != nil {
		panic(err)
	}
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	// Caller can use "github.com/deadblue/gostream/observe" to monitor the
	// uploading progress, please see the example code in that package.
	metadata, err := agent.Upload("0", info, file)
	if err != nil {
		panic(err)
	} else {
		log.Printf("Uploaded file: %#v", metadata)
	}
*/
func (a *Agent) Upload(parentId string, info UploadInfo, r io.Reader) (file *File, err error) {
	// Try close r before method returns.
	rc, ok := r.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(r)
	}
	defer quietly.Close(rc)

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
