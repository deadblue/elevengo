package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
)

const (
	apiUploadInit = "https://uplb.115.com/3.0/sampleinitupload.php"
)

type UploadTicket struct {
	Values    map[string]string
	FileField string
}

// Create an upload ticket.
func (c *Client) CreateUploadTicket(categoryId string, filename string, filesize int64) (ticket *UploadTicket, err error) {
	// Request upload token
	form := core.NewForm().
		WithString("userid", c.ui.UserId).
		WithString("filename", filename).
		WithInt64("filesize", filesize).
		WithString("target", fmt.Sprintf("U_1_%s", categoryId))
	ir := &internal.UploadInitResult{}
	if err = c.hc.JsonApi(apiUploadInit, nil, form, ir); err != nil {
		return
	}
	// Create upload ticket
	ticket = &UploadTicket{
		Values: map[string]string{
			"OSSAccessKeyId": ir.AccessKeyId,
			"key":            ir.ObjectKey,
			"policy":         ir.Policy,
			"callback":       ir.Callback,
			"signature":      ir.Signature,
			"name":           filename,
		},
		FileField: "file",
	}
	return
}

//func (c *Client) upload(categoryId string, storeName string, size int64, data io.Reader) (file *CloudFile, err error) {
//	// get upload parameters
//	form := core.NewForm().
//		WithString("userid", c.info.UserId).
//		WithString("target", fmt.Sprintf("U_1_%s", categoryId)).
//		WithString("filename", storeName).
//		WithInt64("filesize", size)
//	uir := &_FileUploadInitResult{}
//	err = c.requestJson(apiFileUploadInit, nil, form, uir)
//	if err != nil {
//		return
//	}
//	// fill upload form
//	form = core.NewMultipartForm().
//		WithString("OSSAccessKeyId", uir.AccessKeyId).
//		WithString("key", uir.ObjectKey).
//		WithString("policy", uir.Policy).
//		WithString("callback", uir.Callback).
//		WithString("signature", uir.Signature).
//		WithString("name", storeName).
//		WithFile("file", storeName, data)
//	ur := &_FileUploadResult{}
//	err = c.requestJson(uir.UploadUrl, nil, form, ur)
//	if err == nil && !ur.State {
//		err = apiError(ur.Code)
//	}
//	if ur.Data == nil {
//		err = ErrUnexpected
//	} else {
//		file = &CloudFile{
//			IsSystem:   false,
//			IsCategory: false,
//			CategoryId: ur.Data.CategoryId,
//			FileId:     ur.Data.FileId,
//			Name:       ur.Data.FileName,
//			PickCode:   ur.Data.PickCode,
//			Sha1:       ur.Data.Sha1,
//		}
//		file.Size, _ = strconv.ParseInt(ur.Data.FizeSize, 10, 64)
//	}
//	return
//}
