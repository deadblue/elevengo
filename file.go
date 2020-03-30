package elevengo

import (
	"github.com/deadblue/elevengo/core"
	"strconv"
)

type FileListOpts struct {
	offset, limit uint
	sortFlag      string
	sortAsc       bool
}

func (o *FileListOpts) SortByTime(asc bool) *FileListOpts {
	o.sortFlag, o.sortAsc = "user_ptime", asc
	return o
}
func (o *FileListOpts) SortByName(asc bool) *FileListOpts {
	o.sortFlag, o.sortAsc = "file_name", asc
	return o
}
func (o *FileListOpts) SortBySize(asc bool) *FileListOpts {
	o.sortFlag, o.sortAsc = "file_size", asc
	return o
}
func (o *FileListOpts) Next() *FileListOpts {
	o.offset += o.limit
	return o
}
func NewFileListOpts(pageSize uint) *FileListOpts {
	return &FileListOpts{
		offset: 0,
		limit:  pageSize,
	}
}

// Get file list under the specific category.
//
// "0" is a special categoryId which means the root,
// everything starts from here.
//
// `sort` is optional, pass `nil` will use the default sort option:
// sorting by modify time in desc.
func (c *Client) FileList(categoryId string, opts *FileListOpts) (files []*CloudFile, remain int, err error) {
	if opts == nil {
		opts = NewFileListOpts(100).SortByTime(false)
	}
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId).
		WithString("show_dir", "1").
		WithString("snap", "0").
		WithString("natsort", "1").
		WithString("format", "json").
		WithInt("offset", int(opts.offset)).
		WithInt("limit", int(opts.offset)).
		WithString("o", opts.sortFlag)
	if opts.sortAsc {
		qs.WithString("asc", "1")
	} else {
		qs.WithString("asc", "0")
	}
	// call API
	result := &_FileListResult{}
	err = c.requestJson(apiFileList, qs, nil, result)
	if err == nil && !result.State {
		err = apiError(result.MessageCode)
	}
	if err != nil {
		return
	}
	// remain file count
	remain = result.TotalCount - (result.Offset + result.PageSize)
	if remain < 0 {
		remain = 0
	}
	// convert result
	files = make([]*CloudFile, len(result.Data))
	for index, data := range result.Data {
		info := &CloudFile{
			IsCategory: false,
			IsSystem:   (index + result.Offset) < result.SysCount,
			CategoryId: data.CategoryId,
			Name:       data.Name,
			Size:       data.Size,
			PickCode:   data.PickCode,
		}
		info.CreateTime, _ = strconv.ParseInt(data.CreateTime, 10, 64)
		info.UpdateTime, _ = strconv.ParseInt(data.UpdateTime, 10, 64)
		if data.FileId == nil {
			info.IsCategory = true
			info.ParentId = *data.ParentId
		} else {
			info.FileId = *data.FileId
			info.Sha1 = *data.Sha1
		}
		files[index] = info
	}
	return
}

// Search files which's name contains the specific keyword,
// the searching is recursive, starts from the specific category.
//
// `keyword` can not be empty
//
// `offset` is base on zero.
//
// `limit` can not be lower than `FileListMinLimit`,
//  and can not be higher than `FileListMaxLimit`
func (c *Client) FileSearch(categoryId, keyword string, offset, limit int) (files []*CloudFile, remain int, err error) {
	if len(keyword) == 0 {
		return nil, 0, ErrEmptyKeyword
	}
	if limit < FileListMinLimit {
		limit = FileListMinLimit
	} else if limit > FileListMaxLimit {
		limit = FileListMaxLimit
	}
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId).
		WithString("search_value", keyword).
		WithString("format", "json").
		WithInt("offset", offset).
		WithInt("limit", limit)
	// call API
	result := &_FileSearchResult{}
	err = c.requestJson(apiFileSearch, qs, nil, result)
	if err == nil && !result.State {
		err = apiError(result.MessageCode)
	}
	if err != nil {
		return
	}
	// remain file count
	remain = result.TotalCount - (result.Offset + result.PageSize)
	if remain < 0 {
		remain = 0
	}
	// convert result
	files = make([]*CloudFile, len(result.Data))
	for index, data := range result.Data {
		info := &CloudFile{
			IsCategory: false,
			IsSystem:   false,
			CategoryId: data.CategoryId,
			Name:       data.Name,
			Size:       data.Size,
			PickCode:   data.PickCode,
		}
		info.CreateTime, _ = strconv.ParseInt(data.CreateTime, 10, 64)
		info.UpdateTime, _ = strconv.ParseInt(data.UpdateTime, 10, 64)
		if data.FileId == nil {
			info.IsCategory = true
			info.ParentId = *data.ParentId
		} else {
			info.FileId = *data.FileId
			info.Sha1 = *data.Sha1
		}
		files[index] = info
	}
	return
}

func (c *Client) FileRename(fileId, name string) (err error) {
	form := core.NewForm().
		WithString("fid", fileId).
		WithString("file_name", name)
	result := new(_FileOperateResult)
	err = c.requestJson(apiFileEdit, nil, form, result)
	if err == nil && !result.State {
		err = apiError(-1)
	}
	return
}

func (c *Client) FileCopy(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := new(_FileOperateResult)
	err = c.requestJson(apiFileCopy, nil, form, result)
	if err == nil && !result.State {
		err = apiError(-1)
	}
	return
}

func (c *Client) FileMove(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := new(_FileOperateResult)
	err = c.requestJson(apiFileMove, nil, form, result)
	if err == nil && !result.State {
		err = apiError(-1)
	}
	return
}

func (c *Client) FileDelete(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := new(_FileOperateResult)
	err = c.requestJson(apiFileDelete, nil, form, result)
	if err == nil && !result.State {
		err = apiError(-1)
	}
	return
}

func (c *Client) CategoryAdd(parentId, name string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithString("cname", name)
	result := &_FileAddResult{}
	err = c.requestJson(apiFileAdd, nil, form, result)
	if err == nil && !result.State {
		err = apiError(-1)
	}
	return
}

func (c *Client) CategoryInfo(categoryId string) (err error) {
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId)
	result := &CategoryInfoResult{}
	return c.requestJson(apiCategoryGet, qs, nil, result)
}
