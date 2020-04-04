package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"os"
)

const (
	apiUploadInit = "https://uplb.115.com/3.0/sampleinitupload.php"
)

// "UploadInfo" contains all required information to create an upload ticket.
// You need not to implement it, a well-known implementation is "os.FileInfo".
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

// Create an upload ticket.
func (a *Agent) CreateUploadTicket(parentId string, info UploadInfo) (ticket *UploadTicket, err error) {
	// Request upload token
	form := core.NewForm().
		WithInt("userid", a.ui.UserId).
		WithString("filename", info.Name()).
		WithInt64("filesize", info.Size()).
		WithString("target", fmt.Sprintf("U_1_%s", parentId))
	result := &internal.UploadInitResult{}
	if err = a.hc.JsonApi(apiUploadInit, nil, form, result); err != nil {
		return
	}
	// Create upload ticket
	ticket = &UploadTicket{
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

// A simple upload implementation without progress echo.
// I do not suggest using it to upload a large file.
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
	result := &internal.UploadResult{}
	return a.hc.JsonApi(ticket.Endpoint, nil, form, result)
}
