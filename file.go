package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"time"
)

// File describe a file or directory on cloud storage.
type File struct {

	// Marks is the file a directory.
	IsDirectory bool
	// Unique identifier of the file on the cloud storage.
	FileId string
	// FileId of the parent directory.
	ParentId string

	// Base name of the file.
	Name string
	// Size in bytes of the file.
	Size int64
	// Identifier used for downloading or playing the file.
	PickCode string
	// SHA1 hash of file content, in HEX format.
	Sha1 string

	// Is file stared
	Star bool
	// File labels
	Labels []*Label

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

	f.Star = info.IsStar != 0
	f.Labels = make([]*Label, len(info.Labels))
	for i, l := range info.Labels {
		f.Labels[i] = &Label{
			Id:    l.Id,
			Name:  l.Name,
			Color: LabelColor(webapi.LabelColorMap[l.Color]),
		}
	}

	f.CreateTime = time.Unix(int64(info.CreateTime), 0)
	f.UpdateTime = time.Unix(int64(info.UpdateTime), 0)

	return f
}

// DirInfo only used in FileInfo.
type DirInfo struct {
	// Directory ID.
	Id string
	// Directory Name.
	Name string
}

// FileInfo is returned by FileStat(), contains basic information of a file.
type FileInfo struct {

	// Base name of the file.
	Name string
	// Identifier used for downloading or playing the file.
	PickCode string
	// SHA1 hash of file content, in HEX format.
	Sha1 string
	// Marks is file a directory.
	IsDirectory bool
	// Files count under this directory.
	FileCount int
	// Subdirectories count under this directory.
	DirCount int

	// Create time of the file.
	CreateTime time.Time
	// Last update time of the file.
	UpdateTime time.Time
	// Last access time of the file.
	AccessTime time.Time

	// Parent directory list.
	Parents []*DirInfo
}

func fileParseListResponse(resp *webapi.FileListResponse, files []*File, cursor *FileCursor) (n int, err error) {
	// Parse response data
	n = len(files)
	result := make([]*webapi.FileInfo, 0, n)
	if err = resp.Decode(&result); err != nil {
		return 0, err
	}
	// Fill files
	if rn := len(result); rn < n {
		n = rn
	}
	for i := 0; i < n; i++ {
		files[i] = (&File{}).from(result[i])
	}
	// Update cursor
	cursor.total = resp.Count
	cursor.offset += n
	return n, nil
}

// FileList lists files list under a directory whose id is parentId.
func (a *Agent) FileList(dirId string, cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	// Check cursor
	if cursor == nil {
		return 0, webapi.ErrInvalidCursor
	}
	tx := fmt.Sprintf("file_list_%s", dirId)
	if err = cursor.checkTransaction(tx); err != nil {
		return
	}
	// Prepare request
	qs := web.Params{}.
		With("aid", "1").
		With("show_dir", "1").
		With("snap", "0").
		With("natsort", "1").
		With("fc_mix", "0").
		With("format", "json").
		With("cid", dirId).
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
		err = a.wc.CallJsonApi(apiUrl, qs, nil, resp)
		if err == webapi.ErrOrderNotSupport {
			cursor.order, cursor.asc = resp.Order, resp.IsAsc
			qs.With("o", cursor.order).
				WithInt("asc", cursor.asc)
			retry = true
		} else {
			retry = false
		}
	}
	if err != nil {
		return
	}
	// When dirId is not exists, 115 will return the file list under root dir,
	// that should be considered as an error.
	if dirId != string(resp.CategoryId) {
		return 0, webapi.ErrNotExist
	}
	return fileParseListResponse(resp, files, cursor)
}

// FileSearch recursively searches files, whose name contains the keyword and under the directory.
func (a *Agent) FileSearch(dirId, keyword string, cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	// Check cursor
	if cursor == nil {
		return 0, webapi.ErrInvalidCursor
	}
	tx := fmt.Sprintf("file_search_%s_%s", dirId, keyword)
	if err = cursor.checkTransaction(tx); err != nil {
		return
	}
	// Prepare request
	qs := web.Params{}.
		With("aid", "1").
		With("cid", dirId).
		With("search_value", keyword).
		WithInt("offset", cursor.offset).
		WithInt("limit", n).
		With("format", "json")
	resp := &webapi.FileListResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileSearch, qs, nil, resp); err != nil {
		return
	}
	return fileParseListResponse(resp, files, cursor)
}

// FileGet gets file information by its ID.
func (a *Agent) FileGet(fileId string, file *File) (err error) {
	qs := web.Params{}.
		With("file_id", fileId)
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileInfo, qs, nil, resp); err != nil {
		return
	}
	data := &webapi.FileInfo{}
	if err = resp.Decode(data); err == nil {
		file.from(data)
	}
	return
}

// FileStat gets information of a file/directory.
func (a *Agent) FileStat(fileId string, info *FileInfo) (err error) {
	qs := (web.Params{}).With("cid", fileId)
	resp := &webapi.FileStatResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileStat, qs, nil, resp); err != nil {
		return
	}
	info.Name = resp.FileName
	info.PickCode = resp.PickCode
	info.Sha1 = resp.Sha1
	info.CreateTime = time.Unix(int64(resp.CreateTime), 0)
	info.UpdateTime = time.Unix(int64(resp.UpdateTime), 0)
	info.AccessTime = time.Unix(resp.AccessTime, 0)
	// Fill parents
	info.Parents = make([]*DirInfo, len(resp.Paths))
	for i, path := range resp.Paths {
		info.Parents[i] = &DirInfo{
			Id:   string(path.FileId),
			Name: path.FileName,
		}
	}
	// Directory info
	info.IsDirectory = resp.IsFile == 0
	if info.IsDirectory {
		info.FileCount = int(resp.FileCount)
		info.DirCount = int(resp.DirCount)
	}
	return
}

// FileMove moves files into target directory whose id is dirId.
func (a *Agent) FileMove(dirId string, fileIds ...string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	form := web.Params{}.
		With("pid", dirId).
		WithArray("fid", fileIds).
		ToForm()
	return a.wc.CallJsonApi(
		webapi.ApiFileMove, nil, form, &webapi.BasicResponse{})
}

// FileCopy copies files into target directory whose id is dirId.
func (a *Agent) FileCopy(dirId string, fileIds ...string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	form := web.Params{}.
		With("pid", dirId).
		WithArray("fid", fileIds).
		ToForm()
	return a.wc.CallJsonApi(
		webapi.ApiFileCopy, nil, form, &webapi.BasicResponse{})
}

// FileRename renames file to new name.
func (a *Agent) FileRename(fileId, newName string) (err error) {
	form := web.Params{}.
		WithMap("files_new_name", map[string]string{
			fileId: newName,
		}).ToForm()
	return a.wc.CallJsonApi(
		webapi.ApiFileRename, nil, form, &webapi.BasicResponse{})
}

// FileDelete deletes files.
func (a *Agent) FileDelete(fileIds ...string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	form := web.Params{}.WithArray("fid", fileIds).ToForm()
	return a.wc.CallJsonApi(
		webapi.ApiFileDelete, nil, form, &webapi.BasicResponse{})
}

// FileFindDuplications finds all duplicate files which have the same SHA1 hash
// with the given file, and return their ID.
func (a *Agent) FileFindDuplications(fileId string) (dupIds []string, err error) {
	qs := web.Params{}.With("file_id", fileId)
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileFindDuplicate, qs, nil, resp); err != nil {
		return
	}
	// Parse response
	var duplications []*webapi.FileDuplication
	if err = resp.Decode(&duplications); err == nil {
		if size := len(duplications); size > 0 {
			dupIds = make([]string, size)
			for i, dup := range duplications {
				dupIds[i] = dup.FileId
			}
		}
	}
	return
}
