package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"io"
)

const (
	apiUploadInit = "https://uplb.115.com/3.0/sampleinitupload.php"
)

type UploadFile struct {
	Name string
	Size int64
	Data io.Reader
}

func (c *Client) UploadFile(categoryId string, file *UploadFile) (err error) {
	// request upload token
	form := core.NewForm().
		WithString("userid", c.ui.UserId).
		WithString("filename", file.Name).
		WithInt64("filesize", file.Size).
		WithString("target", fmt.Sprintf("U_1_%s", categoryId))
	ir := &internal.UploadInitResult{}
	err = c.hc.JsonApi(apiUploadInit, nil, form, ir)
	// Post data
	form = core.NewMultipartForm()
	form.WithString("OSSAccessKeyId", ir.AccessKeyId).
		WithString("key", ir.ObjectKey).
		WithString("policy", ir.Policy).
		WithString("callback", ir.Callback).
		WithString("signature", ir.Signature).
		WithString("name", file.Name).
		WithFile("file", file.Name, file.Data)
	ur := &internal.UploadResult{}
	return c.hc.JsonApi(ir.Host, nil, form, ur)
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
