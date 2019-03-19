package elevengo

import (
	"fmt"
	"strings"
	"time"
)

func (c *Client) FileList(dirId string, sort OrderMode) (result *FileListResult, err error) {
	params := newRequestParameters().
		With("aid", "1").
		With("cid", dirId).
		With("o", string(sort)).
		With("asc", "1").
		With("show_dir", "1").
		With("snap", "0").
		With("natsort", "1").
		With("format", "json").
		WithInt("offset", 0).
		WithInt("limit", 100)
	result = &FileListResult{}
	err = c.requestJson(apiFileList, params, nil, result)
	return
}

func (c *Client) FileMkdir(parentId, name string) (err error) {
	qs := newRequestParameters().
		With("pid", parentId).
		With("cname", name)
	body := strings.NewReader(qs.QueryString())
	err = c.requestJson(apiFileAdd, nil, body, nil)
	return
}

func (c *Client) FileMove(parentId string, fileIds ...string) (err error) {
	params := newRequestParameters().
		With("pid", parentId).
		WithStrings("fid", fileIds...)
	return c.requestJson(apiFileMove, nil, params.FormData(), nil)
}

func (c *Client) FileRename(fileId, name string) (err error) {
	params := newRequestParameters().
		With("fid", fileId).
		With("file_name", name)
	return c.requestJson(apiFileEdit, nil, params.FormData(), nil)
}

func (c *Client) FileDelete(parentId string, fileIds ...string) {
	qs := newRequestParameters().With("pid", parentId)
	for index, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", index)
		qs.With(key, fileId)
	}
}

func (c *Client) FileSearch(dirId, keyword string) (err error) {
	qs := newRequestParameters().
		With("aid", "1").
		With("cid", dirId).
		With("search_value", keyword).
		With("format", "json").
		WithInt("offset", 0).
		WithInt("limit", 100)
	err = c.requestJson(apiFileSearch, qs, nil, nil)
	return
}

func (c *Client) FileDownload(pickcode string) (result *FileDownloadResult, err error) {
	qs := newRequestParameters().
		With("pickcode", pickcode).
		WithInt64("_", time.Now().UnixNano())
	result = &FileDownloadResult{}
	err = c.requestJson(apiFileDownload, qs, nil, result)
	return
}
