package elevengo

import (
	"errors"
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"time"
)

const (
	apiFileIndex      = "https://webapi.115.com/files/index_info"
	apiFileList       = "https://webapi.115.com/files"
	apiFileListByName = "https://aps.115.com/natsort/files.php"
	apiFileSearch     = "https://webapi.115.com/files/search"
	apiFileAdd        = "https://webapi.115.com/files/add"
	apiFileCopy       = "https://webapi.115.com/files/copy"
	apiFileMove       = "https://webapi.115.com/files/move"
	apiFileRename     = "https://webapi.115.com/files/batch_rename"
	apiFileDelete     = "https://webapi.115.com/rb/delete"

	filePageSizeMin     = 10
	filePageSizeMax     = 1000
	filePageSizeDefault = 100
)

// File page parameter.
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

// Sort file parameter.
type FileSortParam struct {
	flag string
	asc  int
}

// Sort files by update time
func (p *FileSortParam) ByTime() *FileSortParam {
	p.flag = "user_ptime"
	return p
}

// Sort files by name.
func (p *FileSortParam) ByName() *FileSortParam {
	p.flag = "file_name"
	return p
}

// Sort files by size.
func (p *FileSortParam) BySize() *FileSortParam {
	p.flag = "file_size"
	return p
}

// Use ascending order.
func (p *FileSortParam) Asc() *FileSortParam {
	p.asc = 1
	return p
}

// Use descending order.
func (p *FileSortParam) Desc() *FileSortParam {
	p.asc = 0
	return p
}

// CloudFile describe a remote file/category.
// TODO: rename CloudFile to Category.
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

// TODO: Plan to rename this method to "StorageInfo()".
func (a *Agent) FileIndex() (err error) {
	result := new(internal.FileIndexResult)
	err = a.hc.JsonApi(apiFileIndex, nil, nil, result)
	if err != nil {
		// TODO: handle api result
	}
	return
}

// Get one page of files under specific category(directory).
// The remote API can get at most 1000 files in one page, so if
// there are more than 1000 files in a category, you should call
// this API more than 1 times.
func (a *Agent) FileList(parentId string, page *FilePageParam, sort *FileSortParam) (files []*CloudFile, err error) {
	// Prepare parameters
	if sort == nil {
		sort = (&FileSortParam{}).ByTime().Desc()
	}
	if sort.flag == "" {
		sort.flag = "user_ptime"
	}
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", parentId).
		WithString("show_dir", "1").
		WithString("snap", "0").
		WithString("natsort", "1").
		WithString("format", "json").
		WithString("o", sort.flag).
		WithInt("asc", sort.asc).
		WithInt("offset", page.offset()).
		WithInt("limit", page.limit())
	// Select API URL
	apiUrl := apiFileList
	if sort.flag == "file_name" {
		apiUrl = apiFileListByName
	}
	// Call API
	result := &internal.FileListResult{}
	err = a.hc.JsonApi(apiUrl, qs, nil, result)
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
			CreateTime: internal.ParseUnixTime(data.CreateTime),
			UpdateTime: internal.ParseUnixTime(data.UpdateTime),
		}
		if data.FileId != "" {
			f.IsCategory = false
			f.FileId = data.FileId
			f.Size = int64(data.Size)
			f.Sha1 = data.Sha1
		} else {
			f.IsCategory = true
			f.ParentId = data.ParentId
		}
		files[i] = f
	}
	return
}

func (a *Agent) FileSearch(parentId, keyword string, page *FilePageParam) (files []*CloudFile, next bool, err error) {
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", parentId).
		WithString("search_value", keyword).
		WithInt("offset", page.offset()).
		WithInt("limit", page.limit()).
		WithString("format", "json")
	result := &internal.FileSearchResult{}
	err = a.hc.JsonApi(apiFileSearch, qs, nil, result)
	if err != nil {
		return
	}
	return
}

func (a *Agent) FileCopy(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileCopy, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (a *Agent) FileMove(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileMove, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (a *Agent) FileRename(fileId, name string) (err error) {
	form := core.NewForm().
		WithStringMap("files_new_name", map[string]string{fileId: name})
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileRename, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (a *Agent) FileDelete(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileDelete, nil, form, result)
	if err == nil && !result.State {
		// TODO: convert upstream error
		err = errors.New(result.Error)
	}
	return
}

func (a *Agent) CategoryAdd(parentId, name string) (categoryId string, err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithString("cname", name)
	result := &internal.CategoryAddResult{}
	err = a.hc.JsonApi(apiFileAdd, nil, form, result)
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
