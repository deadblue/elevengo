package elevengo

import (
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"time"
)

const (
	apiFileCopy   = "https://webapi.115.com/files/copy"
	apiFileMove   = "https://webapi.115.com/files/move"
	apiFileRename = "https://webapi.115.com/files/batch_rename"
	apiFileDelete = "https://webapi.115.com/rb/delete"

	fileDefaultLimit = 115
)

// File describe a file or directory on cloud storage.
type File struct {

	// IsDirectory marks is the file a directory.
	IsDirectory bool

	// FileId is the unique identifier on the cloud storage.
	FileId string

	// ParentID is the FileId of the parent directory.
	ParentId string

	// Name is name of the file.
	Name string

	// Size is bytes size of the file.
	Size int64

	// PickCode is used for downloading or playing the file.
	PickCode string

	// Sha1 is SHA1 hash value of the file, in HEX format.
	Sha1 string

	// Create time of the file.
	CreateTime time.Time

	// Update time of the file.
	UpdateTime time.Time
}

func (f *File) from(info *webapi.FileInfo) *File {
	if info.FileId != "" {
		f.FileId = info.FileId
		f.ParentId = info.CategoryId
		f.IsDirectory = false
	} else {
		f.FileId = info.CategoryId
		f.ParentId = info.ParentId
		f.IsDirectory = true
	}
	f.Name = info.Name
	f.Size = int64(info.Size)
	f.PickCode = info.PickCode
	f.Sha1 = info.Sha1
	f.CreateTime = time.Unix(int64(info.CreateTime), 0)
	f.UpdateTime = time.Unix(int64(info.UpdateTime), 0)

	return f
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

	// PickCode for downloading or playing.
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

type FileCursor struct {
	init   bool
	offset int
	total  int
	order  string
	asc    int
}

func (c *FileCursor) HasMore() bool {
	return !c.init || c.offset < c.total
}
func (c *FileCursor) Total() int {
	return c.total
}
func (c *FileCursor) Remain() int {
	return c.total - c.offset
}

// FileList lists files list under a directory whose id is parentId.
func (a *Agent) FileList(parentId string, cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	// Initialize cursor
	if !cursor.init {
		cursor.order = "user_ptime"
		cursor.asc = 0
		cursor.init = true
	}
	// Prepare request
	qs := web.Params{}.
		With("aid", "1").
		With("show_dir", "1").
		With("snap", "0").
		With("natsort", "1").
		With("fc_mix", "1").
		With("format", "json").
		With("cid", parentId).
		With("o", cursor.order).
		WithInt("asc", cursor.asc).
		WithInt("offset", cursor.offset).
		WithInt("limit", n)
	resp := &webapi.FileListResponse{}
	for retry := true; retry; {
		// Select API URL
		apiUrl := webapi.ApiFileList
		if cursor.order == "file_name" {
			apiUrl = webapi.ApiFileListByName
		}
		// Call API
		err, retry = a.wc.CallJsonApi(apiUrl, qs, nil, resp), false
		if err != nil {
			break
		}
		// Parse response
		if err = resp.Err(); err != nil {
			if resp.ErrorCode2 == 20130827 {
				// Change order and retry
				cursor.order, cursor.asc = resp.Order, resp.IsAsc
				qs.With("o", cursor.order).WithInt("asc", cursor.asc)
				retry = true
			}
		}
	}
	if err != nil {
		return
	}
	// Upstream will return file list under root when parentId is invalid, but this API should
	// return an error.
	//if parentId != string(resp.CategoryId) {
	//	return 0, errFileNotExist
	//}
	result := make([]*webapi.FileInfo, 0, n)
	if err = resp.Decode(&result); err != nil {
		return
	}
	if rn := len(result); rn < n {
		n = rn
	}
	for i := 0; i < n; i++ {
		files[i] = (&File{}).from(result[i])
	}
	// Update cursor
	cursor.offset += n
	cursor.total = resp.Count
	return
}

// FileSearch recursively searches files, whose name contains the keyword and under the directory.
func (a *Agent) FileSearch(rootId, keyword string, cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	// Initialize cursor
	if !cursor.init {
		cursor.offset = 0
		cursor.init = true
	}
	// Prepare request
	qs := web.Params{}.
		With("aid", "1").
		With("cid", rootId).
		With("search_value", keyword).
		WithInt("offset", cursor.offset).
		WithInt("limit", n).
		With("format", "json")
	resp := webapi.FileSearchResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileSearch, qs, nil, &resp); err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	// Parse response
	result := make([]*webapi.FileInfo, 0, n)
	if err = resp.Decode(&result); err != nil {
		return
	}
	// Fill to files
	if rn := len(result); rn < n {
		n = rn
	}
	for i := 0; i < n; i++ {
		files[i] = (&File{}).from(result[i])
	}
	// Update cursor
	cursor.offset += n
	cursor.total = resp.Count
	return
}

// FileStat gets information of a file/directory.
func (a *Agent) FileStat(fileId string, info *FileInfo) (err error) {
	qs := (web.Params{}).With("cid", fileId)
	resp := &webapi.FileStatResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileStat, qs, nil, resp); err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}

	// TODO: Fill info.
	return
}

// FileCopy copies files into target directory.
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

// FileMove moves files into target directory.
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
