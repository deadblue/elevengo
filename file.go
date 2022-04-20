package elevengo

import (
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"github.com/deadblue/elevengo/internal/util"
	"time"
)

const (
	apiFileIndex      = "https://webapi.115.com/files/index_info"
	apiFileList       = "https://webapi.115.com/files"
	apiFileListByName = "https://aps.115.com/natsort/files.php"
	apiFileStat       = "https://webapi.115.com/category/get"
	apiFileSearch     = "https://webapi.115.com/files/search"
	apiFileCopy       = "https://webapi.115.com/files/copy"
	apiFileMove       = "https://webapi.115.com/files/move"
	apiFileRename     = "https://webapi.115.com/files/batch_rename"
	apiFileDelete     = "https://webapi.115.com/rb/delete"

	fileDefaultLimit = 115

	codeListRetry = 20130827
)

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
	CreateTime time.Time
	// Update time of the file.
	UpdateTime time.Time
}

// DirectoryInfo only used in FileInfo.
type DirectoryInfo struct {
	// Directory ID.
	Id string
	// Directory Name.
	Name string
}

// FileInfo is returned by FileStat(), contains basic information of a file.
type FileInfo struct {
	// True means a file.
	IsFile bool
	// True means a directory.
	IsDirectory bool
	// File name.
	Name string
	// Pick code for downloading.
	PickCode string
	// Sha1 hash value of the file, empty for directory.
	Sha1 string
	// Create time of the file.
	CreateTime time.Time
	// Update time of the file.
	UpdateTime time.Time
	// Parent directory list.
	Parents []*DirectoryInfo
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
	result := &types.FileListResult{}
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
				// Update order flag
				fc.order = result.Order
				fc.asc = result.IsAsc
				// Update query string
				qs.WithString("o", fc.order).WithInt("asc", fc.asc)
				// Try to call API again
				retry = true
			} else {
				err = types.MakeFileError(result.ErrorCode, result.Error)
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
			CreateTime: util.ParseUnixTime(data.CreateTime),
			UpdateTime: util.ParseUnixTime(data.UpdateTime),
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
	result := &types.FileSearchResult{}
	err = a.hc.JsonApi(apiFileSearch, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(result.ErrorCode, result.Error)
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
			CreateTime: util.ParseUnixTime(data.CreateTime),
			UpdateTime: util.ParseUnixTime(data.UpdateTime),
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
	result := &types.FileOperateResult{}
	err = a.hc.JsonApi(apiFileCopy, nil, form, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(int(result.ErrorCode), result.Error)
	}
	return
}

// Move files into specific directory.
func (a *Agent) FileMove(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &types.FileOperateResult{}
	err = a.hc.JsonApi(apiFileMove, nil, form, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(int(result.ErrorCode), result.Error)
	}
	return
}

// Rename file.
func (a *Agent) FileRename(fileId, name string) (err error) {
	form := core.NewForm().
		WithString("fid", fileId).
		WithString("file_name", name).
		WithStringMap("files_new_name", map[string]string{fileId: name})
	result := &types.FileOperateResult{}
	err = a.hc.JsonApi(apiFileRename, nil, form, result)
	if err == nil && !result.State {
		err = types.MakeFileError(int(result.ErrorCode), result.Error)
	}
	return
}

// Delete files from a directory.
func (a *Agent) FileDelete(parentId string, fileIds ...string) (err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithStrings("fid", fileIds)
	result := &types.FileOperateResult{}
	err = a.hc.JsonApi(apiFileDelete, nil, form, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(int(result.ErrorCode), result.Error)
	}
	return
}

/*
Get file information related to the file ID.
Since the upstream response is cheap, this method cat not return more information.
*/
func (a *Agent) FileStat(fileId string) (info FileInfo, err error) {
	qs := core.NewQueryString().
		WithString("aid", "1").
		WithString("cid", fileId)
	result := &types.FileStatResult{}
	err = a.hc.JsonApi(apiFileStat, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = errFileNotExist
	}
	if err == nil {
		data := result.Data
		info = FileInfo{
			IsFile:      data.FileType == "1",
			IsDirectory: data.FileType == "0",
			Name:        data.Name,
			Sha1:        data.Sha1,
			PickCode:    data.PickCode,
			CreateTime:  util.ParseUnixTime(data.CreateTime),
			UpdateTime:  util.ParseUnixTime(data.UpdateTime),
		}
		info.Parents = make([]*DirectoryInfo, len(data.Paths))
		for i, p := range data.Paths {
			info.Parents[i] = &DirectoryInfo{
				Id:   string(p.FileId),
				Name: p.FileName,
			}
		}
	}
	return
}
