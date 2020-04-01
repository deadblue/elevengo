package elevengo

import (
	"errors"
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"strconv"
	"time"
)

const (
	apiFileIndex      = "https://webapi.115.com/files/index_info"
	apiFileList       = "https://webapi.115.com/files"
	apiFileListByName = "https://aps.115.com/natsort/files.php"
	apiFileAdd        = "https://webapi.115.com/files/add"
	apiFileCopy       = "https://webapi.115.com/files/copy"
	apiFileMove       = "https://webapi.115.com/files/move"
	apiFileRename     = "https://webapi.115.com/files/batch_rename"
	apiFileDelete     = "https://webapi.115.com/rb/delete"
	apiFileSearch     = "https://webapi.115.com/files/search"

	filePageSizeMin     = 10
	filePageSizeMax     = 1000
	filePageSizeDefault = 100
)

// Page parameter for `Client.FileList()` and `Client.FileSearch()`
type FilePageParam struct {
	index, size int
}

func (p *FilePageParam) Size(size int) *FilePageParam {
	if size < filePageSizeMin {
		p.size = filePageSizeMin
	} else if size > filePageSizeMax {
		p.size = filePageSizeMax
	} else {
		p.size = size
	}
	return p
}
func (p *FilePageParam) Next() *FilePageParam {
	p.index += 1
	return p
}
func (p *FilePageParam) Prev() *FilePageParam {
	if p.index > 0 {
		p.index -= 1
	}
	return p
}
func (p *FilePageParam) Goto(num int) *FilePageParam {
	if num > 0 {
		p.index = num - 1
	}
	return p
}
func (p *FilePageParam) limit() int {
	if p.size == 0 {
		p.size = filePageSizeDefault
	}
	return p.size
}
func (p *FilePageParam) offset() int {
	return p.index * p.limit()
}

// Sort parameter for `Client.FileList()`
type FileSortParam struct {
	flag string
	asc  int
}

func (p *FileSortParam) ByTime() *FileSortParam {
	p.flag = "user_ptime"
	return p
}
func (p *FileSortParam) ByName() *FileSortParam {
	p.flag = "file_name"
	return p
}
func (p *FileSortParam) BySize() *FileSortParam {
	p.flag = "file_size"
	return p
}
func (p *FileSortParam) Asc() *FileSortParam {
	p.asc = 1
	return p
}
func (p *FileSortParam) Desc() *FileSortParam {
	p.asc = 0
	return p
}

// CloudFile is a remote file/category object
type CloudFile struct {
	IsCategory bool
	FileId     string
	CategoryId string
	ParentId   string
	Name       string
	Size       int64
	PickCode   string
	Sha1       string
	CreateTime time.Time
	UpdateTime time.Time
}

func (c *Client) FileIndex() (err error) {
	result := new(internal.FileIndexResult)
	err = c.hc.JsonApi(apiFileIndex, nil, nil, result)
	if err != nil {
		// TODO: handle api result
	}
	return
}

func parseTimestamp(s string) time.Time {
	sec, _ := strconv.ParseInt(s, 10, 64)
	return time.Unix(sec, 0)
}

func (c *Client) FileList(categoryId string, page *FilePageParam, sort *FileSortParam) (files []*CloudFile, err error) {
	// prepare parameters
	if sort == nil {
		sort = (&FileSortParam{}).ByTime().Desc()
	}
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", categoryId).
		WithString("show_dir", "1").
		WithString("snap", "0").
		WithString("natsort", "1").
		WithString("format", "json").
		WithString("o", sort.flag).
		WithInt("asc", sort.asc).
		WithInt("offset", page.offset()).
		WithInt("limit", page.limit())
	// select API URL
	apiUrl := apiFileList
	if sort.flag == "file_name" {
		apiUrl = apiFileListByName
	}
	// call API
	result := &internal.FileListResult{}
	err = c.hc.JsonApi(apiUrl, qs, nil, result)
	if err == nil && !result.State {
		err = fmt.Errorf("get file list failed")
	}
	if err != nil {
		return
	}
	// Fill files array
	files = make([]*CloudFile, len(result.Data))
	for i, data := range result.Data {
		f := &CloudFile{
			CategoryId: data.CategoryId,
			Name:       data.Name,
			PickCode:   data.PickCode,
			CreateTime: parseTimestamp(data.CreateTime),
			UpdateTime: parseTimestamp(data.UpdateTime),
		}
		if data.FileId != "" {
			f.IsCategory = false
			f.FileId = data.FileId
			f.Size = data.Size
			f.Sha1 = data.Sha1
		} else {
			f.IsCategory = true
			f.ParentId = data.ParentId
		}
		files[i] = f
	}
	return
}

func (c *Client) FileCopy(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = c.hc.JsonApi(apiFileCopy, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (c Client) FileMove(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = c.hc.JsonApi(apiFileMove, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (c *Client) FileRename(fileId, name string) (err error) {
	form := core.NewForm().
		WithStringMap("files_new_name", map[string]string{fileId: name})
	result := &internal.FileOperateResult{}
	err = c.hc.JsonApi(apiFileRename, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (c *Client) FileDelete(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = c.hc.JsonApi(apiFileDelete, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (c *Client) CategoryAdd(parentId, name string) (categoryId string, err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithString("cname", name)
	result := &internal.CategoryAddResult{}
	err = c.hc.JsonApi(apiFileAdd, nil, form, result)
	if err != nil {
		return
	}
	if !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	} else {
		categoryId = result.CategoryId
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
//func (c *Client) FileSearch(categoryId, keyword string, page *PageParam) (files []*CloudFile, remain int, err error) {
//	if len(keyword) == 0 {
//		return nil, 0, ErrEmptyKeyword
//	}
//	if page == nil {
//		page = &PageParam{}
//	}
//	qs := core.NewQueryString().
//		WithString("aid", "1").
//		WithString("cid", categoryId).
//		WithString("search_value", keyword).
//		WithString("format", "json").
//		WithInt("offset", page.offset()).
//		WithInt("limit", page.limit())
//	// call API
//	result := &_FileSearchResult{}
//	err = c.requestJson(apiFileSearch, qs, nil, result)
//	if err == nil && !result.State {
//		err = apiError(result.MessageCode)
//	}
//	if err != nil {
//		return
//	}
//	// remain file count
//	remain = result.TotalCount - (result.Offset + result.PageSize)
//	if remain < 0 {
//		remain = 0
//	}
//	// convert result
//	files = make([]*CloudFile, len(result.Data))
//	for index, data := range result.Data {
//		info := &CloudFile{
//			IsCategory: false,
//			IsSystem:   false,
//			CategoryId: data.CategoryId,
//			Name:       data.Name,
//			Size:       data.Size,
//			PickCode:   data.PickCode,
//		}
//		info.CreateTime, _ = strconv.ParseInt(data.CreateTime, 10, 64)
//		info.UpdateTime, _ = strconv.ParseInt(data.UpdateTime, 10, 64)
//		if data.FileId == nil {
//			info.IsCategory = true
//			info.ParentId = *data.ParentId
//		} else {
//			info.FileId = *data.FileId
//			info.Sha1 = *data.Sha1
//		}
//		files[index] = info
//	}
//	return
//}
