package elevengo

import (
	"bytes"
	"fmt"
	"github.com/deadblue/elevengo/core"
	"io"
	"os"
	"strconv"
)

func (c *Client) UploadFile(categoryId, localFile, storeName string) (file *CloudFile, err error) {
	// open local file
	fp, err := os.OpenFile(localFile, os.O_RDONLY, 0644)
	if err != nil {
		return
	}
	defer fp.Close()
	// get file info
	fi, err := fp.Stat()
	if err != nil {
		return
	}
	if fi.IsDir() {
		return nil, ErrUploadDirectory
	}
	// ready to upload
	if storeName == "" {
		storeName = fi.Name()
	}
	return c.upload(categoryId, storeName, fi.Size(), fp)
}

func (c *Client) UploadData(categoryId, storeName string, data []byte) (file *CloudFile, err error) {
	size, reader := int64(len(data)), bytes.NewReader(data)
	return c.upload(categoryId, storeName, size, reader)
}

func (c *Client) upload(categoryId string, storeName string, size int64, data io.Reader) (file *CloudFile, err error) {
	// get upload parameters
	form := core.NewForm().
		WithString("userid", c.info.UserId).
		WithString("target", fmt.Sprintf("U_1_%s", categoryId)).
		WithString("filename", storeName).
		WithInt64("filesize", size)
	uir := &_FileUploadInitResult{}
	err = c.requestJson(apiFileUploadInit, nil, form, uir)
	if err != nil {
		return
	}
	// fill upload form
	form = core.NewMultipartForm().
		WithString("OSSAccessKeyId", uir.AccessKeyId).
		WithString("key", uir.ObjectKey).
		WithString("policy", uir.Policy).
		WithString("callback", uir.Callback).
		WithString("signature", uir.Signature).
		WithString("name", storeName).
		WithFile("file", storeName, data)
	ur := &_FileUploadResult{}
	err = c.requestJson(uir.UploadUrl, nil, form, ur)
	if err == nil && !ur.State {
		err = apiError(ur.Code)
	}
	if ur.Data == nil {
		err = ErrUnexpected
	} else {
		file = &CloudFile{
			IsSystem:   false,
			IsCategory: false,
			CategoryId: ur.Data.CategoryId,
			FileId:     ur.Data.FileId,
			Name:       ur.Data.FileName,
			PickCode:   ur.Data.PickCode,
			Sha1:       ur.Data.Sha1,
		}
		file.Size, _ = strconv.ParseInt(ur.Data.FizeSize, 10, 64)
	}
	return
}
