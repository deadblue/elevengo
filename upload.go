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

type UploadInfo interface {
	Name() string
	Size() int64
}

type UploadTicket struct {
	Endpoint  string
	Values    map[string]string
	FileField string
}

// Create an upload ticket.
func (c *Client) CreateUploadTicket(parentId string, info UploadInfo) (ticket *UploadTicket, err error) {
	// Request upload token
	form := core.NewForm().
		WithInt("userid", c.ui.UserId).
		WithString("filename", info.Name()).
		WithInt64("filesize", info.Size()).
		WithString("target", fmt.Sprintf("U_1_%s", parentId))
	result := &internal.UploadInitResult{}
	if err = c.hc.JsonApi(apiUploadInit, nil, form, result); err != nil {
		return
	}
	// Create upload ticket
	ticket = &UploadTicket{
		Endpoint: result.Host,
		Values: map[string]string{
			"OSSAccessKeyId": result.AccessKeyId,
			"key":            result.ObjectKey,
			"policy":         result.Policy,
			"callback":       result.Callback,
			"signature":      result.Signature,
			"name":           info.Name(),
		},
		FileField: "file",
	}
	return
}

func (c *Client) UploadFile(parentId, localFile string) (err error) {
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
	ticket, err := c.CreateUploadTicket(parentId, info)
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
	return c.hc.JsonApi(ticket.Endpoint, nil, form, result)
}
