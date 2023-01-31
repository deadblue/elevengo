package elevengo

import (
	"time"

	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
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

	// Deprecated
	// Create time of the file. 
	CreateTime time.Time
	// Deprecated
	// Update time of the file.
	UpdateTime time.Time

	// Last modified time
	ModifiedTime time.Time
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

	if info.UpdatedTime != "" {
		f.ModifiedTime = webapi.ParseFileTime(info.UpdatedTime)
	} else {
		f.ModifiedTime = webapi.ParseFileTime(info.ModifiedTime)
	}

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

type fileIterator struct {
	// Common parameters
	dirId  string
	offset int
	order  string
	asc    int
	// Function parameters
	params map[string]string
	// Total count
	count int
	// Cached files
	files []*webapi.FileInfo
	// Cache index
	index int
	// Cache size
	size int
	// Update function
	update func(*fileIterator) error
}

func (i *fileIterator) Next() (err error) {
	if i.index += 1; i.index < i.size {
		return
	}
	i.offset += i.size
	if i.offset >= i.count {
		return webapi.ErrReachEnd
	}
	return i.update(i)
}

func (i *fileIterator) Index() int {
	return i.offset + i.index
}

func (i *fileIterator) Get(file *File) error {
	if i.index >= i.size {
		return webapi.ErrReachEnd
	}
	if file != nil {
		file.from(i.files[i.index])
	}
	return nil
}

func (i *fileIterator) Count() int {
	return i.count
}

// FileIterate returns an iterator.
func (a *Agent) FileIterate(dirId string) (it Iterator[File], err error) {
	fi := &fileIterator{
		dirId:  dirId,
		order:  webapi.FileOrderByTime,
		asc:    0,
		offset: 0,
		update: a.fileListInternal,
	}
	if err = a.fileListInternal(fi); err == nil {
		it = fi
	}
	return
}

// FileStared lists all stared files.
func (a *Agent) FileStared() (it Iterator[File], err error) {
	fi := &fileIterator{
		dirId:  "0",
		order:  webapi.FileOrderByName,
		asc:    0,
		offset: 0,
		params: map[string]string{
			"star": "1",
		},
		update: a.fileListInternal,
	}
	if err = a.fileListInternal(fi); err == nil {
		it = fi
	}
	return
}

func (a *Agent) fileListInternal(fi *fileIterator) (err error) {
	// Prepare request
	qs := web.Params{}.
		With("aid", "1").
		With("show_dir", "1").
		With("snap", "0").
		With("natsort", "1").
		With("fc_mix", "0").
		With("format", "json").
		With("cid", fi.dirId).
		With("o", fi.order).
		WithInt("asc", fi.asc).
		WithInt("offset", fi.offset).
		WithInt("limit", webapi.FileListLimit)
	for pn, pv := range fi.params {
		qs.With(pn, pv)
	}
	resp := &webapi.FileListResponse{}
	for retry := true; retry; {
		// Select API URL
		apiUrl := webapi.ApiFileList
		if fi.order == webapi.FileOrderByName {
			apiUrl = webapi.ApiFileListByName
		}
		// Call API
		err = a.wc.CallJsonApi(apiUrl, qs, nil, resp)
		if err == webapi.ErrOrderNotSupport {
			// Update order & asc
			fi.order, fi.asc = resp.Order, resp.IsAsc
			qs.With("o", fi.order).WithInt("asc", fi.asc)
			retry = true
		} else {
			retry = false
		}
	}
	if err != nil {
		return
	}
	// When dirId not exists, 115 will return the files under root dir, that
	// should be considered as an error.
	if fi.dirId != string(resp.CategoryId) {
		return webapi.ErrNotExist
	}
	// Parse response
	if fi.count = resp.Count; fi.count > 0 {
		fi.files = make([]*webapi.FileInfo, 0, webapi.FileListLimit)
		if err = resp.Decode(&fi.files); err != nil {
			return
		}
		fi.index, fi.size = 0, len(fi.files)
	}
	return
}

// FileSearch recursively searches files under dirId, whose name contains keyword.
func (a *Agent) FileSearch(dirId, keyword string) (it Iterator[File], err error) {
	fi := &fileIterator{
		dirId:  dirId,
		offset: 0,
		params: map[string]string{
			"search_value": keyword,
		},
		update: a.fileSearchInternal,
	}
	if err = a.fileSearchInternal(fi); err == nil {
		it = fi
	}
	return
}

// FileLabeled lists all files which has specific label.
func (a *Agent) FileLabeled(labelId string) (it Iterator[File], err error) {
	fi := &fileIterator{
		dirId:  "0",
		offset: 0,
		params: map[string]string{
			"file_label": labelId,
		},
		update: a.fileSearchInternal,
	}
	if err = a.fileSearchInternal(fi); err == nil {
		it = fi
	}
	return
}

func (a *Agent) fileSearchInternal(fi *fileIterator) (err error) {
	// Prepare request
	qs := web.Params{}.
		With("aid", "1").
		With("cid", fi.dirId).
		WithInt("offset", fi.offset).
		WithInt("limit", webapi.FileListLimit).
		With("format", "json")
	for pn, pv := range fi.params {
		qs.With(pn, pv)
	}
	resp := &webapi.FileListResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileSearch, qs, nil, resp); err != nil {
		return
	}
	// Parse response
	if fi.count = resp.Count; fi.count > 0 {
		fi.files = make([]*webapi.FileInfo, 0, webapi.FileListLimit)
		if err = resp.Decode(&fi.files); err != nil {
			return
		}
		fi.index, fi.size = 0, len(fi.files)
	}
	return
}

// FileGet gets information of a file/directory by its ID.
func (a *Agent) FileGet(fileId string, file *File) (err error) {
	qs := web.Params{}.
		With("file_id", fileId)
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileInfo, qs, nil, resp); err != nil {
		return
	}
	data := make([]*webapi.FileInfo, 0, 1)
	if err = resp.Decode(&data); err == nil {
		file.from(data[0])
	}
	return
}

// FileStat gets statistic information of a file/directory.
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
