package elevengo

import (
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

	codeListRetry = 20130827
)

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

// Create a file cursor for "Agent.FileList()" and "Agent.FileSearch()".
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
	// File name.
	Name string
	// Sha1 hash value of the file, empty for directory.
	Sha1 string
	// Pick code for downloading.
	PickCode string
	// Create time of the file.
	CreateTime *time.Time
	// Update time of the file.
	UpdateTime *time.Time
	// Parent directory ID list.
	ParentIds []string
}

// Get storage size information.
func (a *Agent) StorageStat() (info *StorageInfo, err error) {
	result := new(internal.FileIndexResult)
	err = a.hc.JsonApi(apiFileIndex, nil, nil, result)
	if err == nil && result.IsFailed() {
		err = internal.MakeFileError(result.Code, result.Error)
	}
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

/*
Get file list from specific directory.

The upstream API restricts the data count in response, so for a directory which contains
a lot of files. you need pass a cursor to receive the cursor information, and use it to
fetch remain files.

The cursor should be created by FileCursor(), and DO NOT pass it as nil even you try to
get file list from a empty directory.
*/
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
	result := &internal.FileListResult{}
	for retry := true; retry; {
		// Select API URL
		apiUrl := apiFileList
		if fc.order == "file_name" {
			apiUrl = apiFileListByName
		}
		// Call API
		err = a.hc.JsonApi(apiUrl, qs, nil, result)
		retry = false
		// Handle error
		if err == nil && result.IsFailed() {
			if result.ErrorCode == codeListRetry {
				// Update query string
				qs.WithString("o", fc.order).WithInt("asc", fc.asc)
				// Update order flag
				fc.order = result.Order
				fc.asc = result.IsAsc
				// Try to call API again
				retry = true
			} else {
				err = internal.MakeFileError(result.ErrorCode, result.Error)
			}
		}
	}
	// Upstream will return file list under root when parentId is invalid, but this API should
	// return an error.
	if parentId != string(result.CategoryId) {
		err = errFileNotExist
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

// Recursively search files which's name contains the keyword and under the directory.
func (a *Agent) FileSearch(rootId, keyword string, cursor Cursor) (files []*File, err error) {
	fc, ok := cursor.(*fileCursor)
	if !ok {
		return nil, errFileCursorInvalid
	}
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", rootId).
		WithString("search_value", keyword).
		WithInt("offset", fc.offset).
		WithInt("limit", fc.limit).
		WithString("format", "json")
	result := &internal.FileSearchResult{}
	err = a.hc.JsonApi(apiFileSearch, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = internal.MakeFileError(result.ErrorCode, result.Error)
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

// Copy files into specific directory.
func (a *Agent) FileCopy(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileCopy, nil, form, result)
	if err == nil && result.IsFailed() {
		err = internal.MakeFileError(int(result.ErrorCode), result.Error)
	}
	return
}

// Move files into specific directory.
func (a *Agent) FileMove(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &internal.FileOperateResult{}
	err = a.hc.JsonApi(apiFileMove, nil, form, result)
	if err == nil && result.IsFailed() {
		err = internal.MakeFileError(int(result.ErrorCode), result.Error)
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
		err = internal.MakeFileError(int(result.ErrorCode), result.Error)
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
		err = internal.MakeFileError(int(result.ErrorCode), result.Error)
	}
	return
}

// Create a directory under a directory with specific name.
func (a *Agent) FileMkdir(parentId, name string) (directoryId string, err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithString("cname", name)
	result := &internal.FileAddResult{}
	err = a.hc.JsonApi(apiFileAdd, nil, form, result)
	if err == nil && result.IsFailed() {
		err = internal.MakeFileError(int(result.ErrorCode), result.Error)
	}
	if err == nil {
		directoryId = result.CategoryId
	}
	return
}

/*
Get file information related to the file ID.
Since the upstream response is cheap, this method cat not return more information.
*/
func (a *Agent) FileStat(fileId string) (info *FileInfo, err error) {
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", fileId)
	result := &internal.FileStatResult{}
	err = a.hc.JsonApi(apiFileStat, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = errFileNotExist
	}
	if err == nil {
		data := result.Data
		info = &FileInfo{
			IsFile:      data.FileType == "1",
			IsDirectory: data.FileType == "0",
			Name:        data.Name,
			Sha1:        data.Sha1,
			PickCode:    data.PickCode,
			CreateTime:  internal.ParseUnixTime(data.CreateTime),
			UpdateTime:  internal.ParseUnixTime(data.UpdateTime),
		}
		info.ParentIds = make([]string, len(data.Paths))
		for i, p := range data.Paths {
			info.ParentIds[i] = string(p.FileId)
		}
	}
	return
}
