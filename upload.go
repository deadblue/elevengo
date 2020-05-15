package elevengo

import (
	"encoding/json"
	"fmt"
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"os"
	"time"
)

const (
	apiUploadInit = "https://uplb.115.com/3.0/sampleinitupload.php"
)

/*
UploadInfo contains all required information to create an upload ticket.

To upload a regular file, caller can use os.FileInfo as UploadInfo, see
"Agent.CreateUploadTicket" doc for detail.

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
Create an upload ticket, caller can use thirdparty libraries/tools to process this ticket.

Example with curl:

	filename := "/path/to/file"
	// Get file info
	info, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	// Create upload ticket
	ticket, err := agent.CreateUploadTicket(parentId, info)
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
	file, err := agent.ParseUploadResult(response)
	if err != nil {
		log.Fatel(f)
	} else {
		log.Printf("Uploaded file: %#v", file)
	}
*/
func (a *Agent) CreateUploadTicket(parentId string, info UploadInfo) (ticket UploadTicket, err error) {
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

// Parse uploading response, see "CreateUploadTicket()" doc for detail.
func (a *Agent) ParseUploadResult(content []byte) (file *File, err error) {
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

// A simple upload implementation without progress echo.
func (a *Agent) UploadFile(parentId, localFile string) (err error) {
	// Open local file
	file, err := os.Open(localFile)
	if err != nil {
		return
	}
	// Get file information (should contains name and size)
	info, err := file.Stat()
	if err != nil {
		return nil
	}
	// Create upload ticket
	ticket, err := a.CreateUploadTicket(parentId, info)
	if err != nil {
		return nil
	}
	// Upload file
	form := core.NewMultipartForm().
		WithFile(ticket.FileField, info.Name(), file)
	for name, value := range ticket.Values {
		form.WithString(name, value)
	}
	result := &types.UploadResult{}
	return a.hc.JsonApi(ticket.Endpoint, nil, form, result)
}
