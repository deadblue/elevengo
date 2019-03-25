package elevengo

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func (c *Client) UploadFile(categoryId, localFile, storeName string) (file *UploadedFile, err error) {
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

func (c *Client) UploadData(categoryId, storeName string, data []byte) (file *UploadedFile, err error) {
	size, reader := int64(len(data)), bytes.NewReader(data)
	return c.upload(categoryId, storeName, size, reader)
}

func (c *Client) upload(categoryId string, storeName string, size int64, data io.Reader) (file *UploadedFile, err error) {
	// get upload parameters
	form := newForm(false).
		WithString("userid", c.info.UserId).
		WithString("target", fmt.Sprintf("U_1_%s", categoryId)).
		WithString("filename", storeName).
		WithInt64("filesize", size)
	ir := &_UploadInitResult{}
	err = c.requestJson(apiUploadInit, nil, form, ir)
	if err != nil {
		return
	}
	// fill upload form
	form = newForm(true).
		WithString("OSSAccessKeyId", ir.AccessKeyId).
		WithString("key", ir.ObjectKey).
		WithString("policy", ir.Policy).
		WithString("callback", ir.Callback).
		WithString("signature", ir.Signature).
		WithString("name", storeName).
		WithFile("file", storeName, data)
	ur := &_UploadResult{}
	err = c.requestJson(ir.UploadUrl, nil, form, ur)
	if err == nil && !ur.State {
		err = apiError(ur.Code)
	}
	return ur.Data, err
}
