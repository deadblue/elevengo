package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"time"
)

const (
	apiFileIndex      = "https://webapi.115.com/files/index_info"
	apiFileList       = "https://webapi.115.com/files"
	apiFileListByName = "https://aps.115.com/natsort/files.php"
	apiFileStat       = "https://webapi.115.com/category/get"
	apiFileSearch     = "https://webapi.115.com/files/search"
	apiFileAdd        = "https://webapi.115.com/files/add"
	apiFileCopy       = "https://webapi.115.com/files/copy"
	apiFileMove       = "https://webapi.115.com/files/move"
	apiFileRename     = "https://webapi.115.com/files/batch_rename"
	apiFileDelete     = "https://webapi.115.com/rb/delete"

	filePageSizeMin     = 10
	filePageSizeMax     = 1150
	filePageSizeDefault = 100

	errDirectoryExisting = 20004
)

type FileError int

func (e FileError) Error() string {
	return fmt.Sprintf("remote error: %d", e)
}

func (e FileError) IsAlreadyExisting() bool {
	return e == errDirectoryExisting
}

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

// Storage information.
type StorageInfo struct {
	// Total size in bytes.
	Size int64
	// Used size in bytes.
	Used int64
	// Avail size in bytes.
	Avail int64
}

// File describe a remote file or directory.
type File struct {
	// True of a file, false for a directory.
	IsFile bool
	// Unique ID for the file.
	FileId string
	// Parent directory ID.
	ParentId string
	// File name.
	Name string
	// File size in bytes, 0 for directory.
	Size int64
	// Pick code, you can use this to create a download ticket.
	PickCode string
	// Sha1 hash value of the file, empty for directory.
	Sha1 string
	// Create time of the file.
	CreateTime *time.Time
	// Update time of the file.
	UpdateTime *time.Time
}

// Get storage size information.
func (a *Agent) StorageStat() (info *StorageInfo, err error) {
	result := new(internal.FileIndexResult)
	err = a.hc.JsonApi(apiFileIndex, nil, nil, result)
	if err != nil {
		return
	}
	info = &StorageInfo{
		Size:  int64(result.Data.SpaceInfo.AllTotal.Size),
		Used:  int64(result.Data.SpaceInfo.AllUsed.Size),
		Avail: int64(result.Data.SpaceInfo.AllRemain.Size),
	}
	return
}

// Get one page of files under specific directory.
// The remote API can get at most 1000 files in one page, so if there are
// more than 1000 files in a category, you should call this method multiple times.
func (a *Agent) FileList(parentId string, page *FilePageParam, sort *FileSortParam) (files []*File, err error) {
	// Prepare parameters
	if sort == nil {
		sort = (&FileSortParam{}).ByTime().Desc()
	} else if sort.flag == "" {
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
	if err == nil && result.IsFailed() {
		err = fmt.Errorf("get file list failed")
	}
	if err != nil {
		return
	}
	// Fill files array
	files = make([]*File, len(result.Data))
	for i, data := range result.Data {
		f := &File{
			Name:       data.Name,
			Size:       int64(data.Size),
			PickCode:   data.PickCode,
			Sha1:       data.Sha1,
			CreateTime: internal.ParseUnixTime(data.CreateTime),
			UpdateTime: internal.ParseUnixTime(data.UpdateTime),
		}
		if data.FileId != "" {
			f.IsFile = true
			f.FileId = data.FileId
			f.ParentId = data.CategoryId
		} else {
			f.IsFile = false
			f.FileId = data.CategoryId
			f.ParentId = data.ParentId
		}
		files[i] = f
	}
	// TODO: Plan to update "page" parameter.
	return
}

func (a *Agent) FileSearch(parentId, keyword string, page *FilePageParam) (files []*File, err error) {
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", parentId).
		WithString("search_value", keyword).
		WithInt("offset", page.offset()).
		WithInt("limit", page.limit()).
		WithString("format", "json")
	result := &internal.FileSearchResult{}
	err = a.hc.JsonApi(apiFileSearch, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = fmt.Errorf("remote API error: %s", result.Error)
	}
	if err != nil {
		return
	}
	// Convert result to "File" slice
	files = make([]*File, len(result.Data))
	for i, data := range result.Data {
		files[i] = &File{
			Name:       data.Name,
			Size:       int64(data.Size),
			PickCode:   data.PickCode,
			Sha1:       data.Sha1,
			CreateTime: internal.ParseUnixTime(data.CreateTime),
			UpdateTime: internal.ParseUnixTime(data.UpdateTime),
		}
		if data.FileId != "" {
			files[i].IsFile = true
			files[i].FileId = data.FileId
			files[i].ParentId = data.CategoryId
		} else {
			files[i].IsFile = false
			files[i].FileId = data.CategoryId
			files[i].ParentId = data.ParentId
		}
	}
	// TODO: Plan to update "page" parameter.
	return
}

// Copy files to a directory.
func (a *Agent) FileCopy(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileCopy, nil, form, result)
	if err == nil && result.IsFailed() {
		err = FileError(result.ErrorCode)
	}
	return
}

// Move files to a directory.
func (a *Agent) FileMove(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileMove, nil, form, result)
	if err == nil && result.IsFailed() {
		err = FileError(result.ErrorCode)
	}
	return
}

// Rename file.
func (a *Agent) FileRename(fileId, name string) (err error) {
	form := core.NewForm().
		WithString("fid", fileId).
		WithString("file_name", name).
		WithStringMap("files_new_name", map[string]string{fileId: name})
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileRename, nil, form, result)
	if err == nil && !result.State {
		err = FileError(result.ErrorCode)
	}
	return
}

// Delete files from a directory.
func (a *Agent) FileDelete(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileDelete, nil, form, result)
	if err == nil && result.IsFailed() {
		err = FileError(result.ErrorCode)
	}
	return
}

// Create a directory under specific parent driectory with specific name.
func (a *Agent) FileMkdir(parentId, name string) (categoryId string, err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithString("cname", name)
	result := &internal.FileAddResult{}
	err = a.hc.JsonApi(apiFileAdd, nil, form, result)
	if err == nil && result.IsFailed() {
		err = FileError(result.ErrorCode)
	}
	if err == nil {
		categoryId = result.CategoryId
	}
	return
}

// Get remote file info.
func (a *Agent) FileStat(fileId string) (err error) {
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", fileId)
	result := &internal.FileStatResult{}
	err = a.hc.JsonApi(apiFileStat, qs, nil, result)
	if err == nil {
		a.l.Info(fmt.Sprintf("Stat: %#v", result))
	}
	// TODO: T.B.D
	return
}
