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

	fileDefaultLimit = 100

	errListRetry = 20130827

	errDirectoryExisting = 20004
)

type FileError int

func (e FileError) Error() string {
	return fmt.Sprintf("remote error: %d", e)
}

func (e FileError) IsAlreadyExisting() bool {
	return e == errDirectoryExisting
}

type fileCursor struct {
	used   bool
	order  string
	asc    int
	offset int
	limit  int
	total  int
}

func (c *fileCursor) HasMore() bool {
	return !c.used || c.offset < c.total
}
func (c *fileCursor) Next() {
	c.offset += c.limit
}
func (c *fileCursor) Total() int {
	return c.total
}
func FileCursor() Cursor {
	return &fileCursor{
		used:   false,
		order:  "user_ptime",
		asc:    0,
		offset: 0,
		limit:  fileDefaultLimit,
		total:  0,
	}
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
	// True means a file.
	IsFile bool
	// True means a directory.
	IsDirectory bool
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

// FileInfo is returned by FileStat(), contains basic information of a file.
type FileInfo struct {
	// True means a file.
	IsFile bool
	// True means a directory.
	IsDirectory bool
	// Path
	ParentIds []string
	// Donwload code
	PickCode string
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

// List files under specific directory.
// TODO: Update the doc.
func (a *Agent) FileList(parentId string, cursor Cursor) (files []*File, err error) {
	fc, ok := cursor.(*fileCursor)
	if !ok {
		return nil, errFileCursorInvalid
	}
	// Prepare parameters
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", parentId).
		WithString("show_dir", "1").
		WithString("snap", "0").
		WithString("natsort", "1").
		WithString("format", "json").
		WithString("o", fc.order).
		WithInt("asc", fc.asc).
		WithInt("offset", fc.offset).
		WithInt("limit", fc.limit)
	retry, result := true, &internal.FileListResult{}
	for retry {
		// Select API URL
		apiUrl := apiFileList
		if fc.order == "file_name" {
			apiUrl = apiFileListByName
		}
		// Call API
		err = a.hc.JsonApi(apiUrl, qs, nil, result)
		// Handle error
		retry = false
		if err == nil && result.IsFailed() {
			if result.ErrorCode == errListRetry {
				// Update query string
				qs.WithString("o", fc.order).WithInt("asc", fc.asc)
				// Update order flag
				fc.order = result.Order
				fc.asc = result.IsAsc
				// Try to call API again
				retry = true
			} else {
				err = FileError(result.ErrorCode)
			}
		}
	}
	if err != nil {
		return
	}
	// Update cursor
	fc.used, fc.total = true, result.Count
	fc.offset, fc.limit = result.Offset, result.PageSize
	// Fill files array
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
			files[i].IsDirectory = false
			files[i].FileId = data.FileId
			files[i].ParentId = data.CategoryId
		} else {
			files[i].IsFile = false
			files[i].IsDirectory = true
			files[i].FileId = data.CategoryId
			files[i].ParentId = data.ParentId
		}
	}
	return
}

// Search files which's name contains the specific keyword and under the specific directory.
// TODO: Update the doc.
func (a *Agent) FileSearch(parentId, keyword string, cursor Cursor) (files []*File, err error) {
	fc, ok := cursor.(*fileCursor)
	if !ok {
		return nil, errFileCursorInvalid
	}
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", parentId).
		WithString("search_value", keyword).
		WithInt("offset", fc.offset).
		WithInt("limit", fc.limit).
		WithString("format", "json")
	result := &internal.FileSearchResult{}
	err = a.hc.JsonApi(apiFileSearch, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = FileError(result.ErrorCode)
	}
	if err != nil {
		return
	}
	// Update cursor
	fc.used, fc.total = true, result.Count
	fc.offset, fc.limit = result.Offset, result.PageSize
	// Fill result files
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
			files[i].IsDirectory = false
			files[i].FileId = data.FileId
			files[i].ParentId = data.CategoryId
		} else {
			files[i].IsFile = false
			files[i].IsDirectory = true
			files[i].FileId = data.CategoryId
			files[i].ParentId = data.ParentId
		}
	}
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

// Create a directory under specific parent directory with specific name.
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
func (a *Agent) FileStat(fileId string) (info *FileInfo, err error) {
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", fileId)
	result := &internal.FileStatResult{}
	err = a.hc.JsonApi(apiFileStat, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = errFileStatFailed
	}
	if err == nil {
		data := result.Data
		info = &FileInfo{
			IsFile:      data.FileType == "1",
			IsDirectory: data.FileType == "0",
			PickCode:    data.PickCode,
		}
		info.ParentIds = make([]string, len(data.Paths))
		for i, p := range data.Paths {
			info.ParentIds[i] = string(p.FileId)
		}
	}
	return
}
