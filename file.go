package elevengo

import (
	"time"

	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/api/errors"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/elevengo/option"
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

	// Last modified time
	ModifiedTime time.Time
}

func (f *File) from(info *api.FileInfo) *File {
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
	f.Size = info.Size.Int64()
	f.PickCode = info.PickCode
	f.Sha1 = info.Sha1

	f.Star = bool(info.IsStar)
	// f.Labels = make([]*Label, len(info.Labels))
	// for i, l := range info.Labels {
	// 	f.Labels[i] = &Label{
	// 		Id:    l.Id,
	// 		Name:  l.Name,
	// 		Color: LabelColor(webapi.LabelColorMap[l.Color]),
	// 	}
	// }

	if info.UpdatedTime != "" {
		f.ModifiedTime = api.ParseFileTime(info.UpdatedTime)
	} else {
		f.ModifiedTime = api.ParseFileTime(info.ModifiedTime)
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
	// Parameters
	dirId   string
	offset  int
	order   string
	asc     int
	options []option.FileOption

	// Total count
	count int

	// Cache index
	index int
	// Cache size
	size int
	// Cached files
	files []*api.FileInfo

	// Update function
	update func(*fileIterator) error
}

func (i *fileIterator) Next() (err error) {
	if i.index += 1; i.index < i.size {
		return
	}
	i.offset += i.size
	if i.offset >= i.count {
		return errors.ErrReachEnd
	}
	return i.update(i)
}

func (i *fileIterator) Index() int {
	return i.offset + i.index
}

func (i *fileIterator) Count() int {
	return i.count
}

func (i *fileIterator) Get(file *File) error {
	if i.index >= i.size {
		return errors.ErrReachEnd
	}
	if file != nil {
		file.from(i.files[i.index])
	}
	return nil
}

// FileIterate returns an iterator.
func (a *Agent) FileIterate(dirId string, options ...option.FileOption) (it Iterator[File], err error) {
	fi := &fileIterator{
		dirId:   dirId,
		offset:  0,
		order:   api.FileOrderDefault,
		asc:     0,
		options: options,
		update:  a.fileListInternal,
	}
	if err = a.fileListInternal(fi); err == nil {
		it = fi
	}
	return
}

// FileStared lists all stared files.
// func (a *Agent) FileStared(options ...option.FileOption) (it Iterator[File], err error) {
// 	fi := &fileIterator{
// 		dirId:  "0",
// 		order:  webapi.FileOrderByName,
// 		asc:    0,
// 		offset: 0,
// 		params: map[string]string{
// 			"star": "1",
// 		},
// 		update: a.fileListInternal,
// 	}
// 	if err = a.fileListInternal(fi); err == nil {
// 		it = fi
// 	}
// 	return
// }

func (a *Agent) fileListInternal(fi *fileIterator) (err error) {
	spec := (&api.FileListSpec{}).Init(fi.dirId, fi.offset)
	spec.SetOrder(fi.order, fi.asc)
	for _, opt := range fi.options {
		switch opt := opt.(type) {
		case option.FileStarOption:
			if opt {
				spec.SetStared()
			}
		}
	}
	for retry := true; retry; {
		if err = a.pc.ExecuteApi(spec); err != nil {
			if ferr, ok := err.(*errors.ErrFileOrderNotSupported); ok {
				spec.SetOrder(ferr.Order, ferr.Asc)
			} else {
				return err
			}
		} else {
			retry = false
		}
	}
	result := spec.Result
	fi.order, fi.asc = result.Order, result.Asc
	if fi.count = result.Count; fi.count > 0 {
		fi.index, fi.size = 0, len(result.Files)
		fi.files = make([]*api.FileInfo, fi.size)
		copy(fi.files, result.Files)
	}
	return
}

// FileSearch recursively searches files under dirId, whose name contains keyword.
func (a *Agent) FileSearch(dirId string, options ...option.FileOption) (it Iterator[File], err error) {
	// TODO
	// fi := &fileIterator{
	// 	dirId:  dirId,
	// 	offset: 0,
	// 	params: map[string]string{
	// 		"search_value": keyword,
	// 	},
	// 	update: a.fileSearchInternal,
	// }
	// fi.applyOptions(options)
	// if err = a.fileSearchInternal(fi); err == nil {
	// 	it = fi
	// }
	return
}

// FileLabeled lists all files which has specific label.
func (a *Agent) FileLabeled(labelId string, options ...option.FileOption) (it Iterator[File], err error) {
	// TODO
	// fi := &fileIterator{
	// 	dirId:  "0",
	// 	offset: 0,
	// 	params: map[string]string{
	// 		"file_label": labelId,
	// 	},
	// 	update: a.fileSearchInternal,
	// }
	// fi.applyOptions(options)
	// if err = a.fileSearchInternal(fi); err == nil {
	// 	it = fi
	// }
	return
}

func (a *Agent) fileSearchInternal(fi *fileIterator) (err error) {
	// TODO
	// Prepare request
	// qs := protocol.Params{}.
	// 	With("aid", "1").
	// 	With("cid", fi.dirId).
	// 	WithInt("offset", fi.offset).
	// 	WithInt("limit", webapi.FileListLimit).
	// 	With("format", "json")
	// for pn, pv := range fi.params {
	// 	qs.With(pn, pv)
	// }
	// resp := &webapi.FileListResponse{}
	// if err = a.pc.CallJsonApi(webapi.ApiFileSearch, qs, nil, resp); err != nil {
	// 	return
	// }
	// // Parse response
	// if fi.count = resp.Count; fi.count > 0 {
	// 	fi.files = make([]*webapi.FileInfo, 0, webapi.FileListLimit)
	// 	if err = resp.Decode(&fi.files); err != nil {
	// 		return
	// 	}
	// 	fi.index, fi.size = 0, len(fi.files)
	// }
	return
}

// FileGet gets information of a file/directory by its ID.
func (a *Agent) FileGet(fileId string, file *File) (err error) {
	qs := protocol.Params{}.
		With("file_id", fileId)
	resp := &webapi.BasicResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiFileInfo, qs, nil, resp); err != nil {
		return
	}
	data := make([]*api.FileInfo, 0, 1)
	if err = resp.Decode(&data); err == nil {
		file.from(data[0])
	}
	return
}

// FileStat gets statistic information of a file/directory.
func (a *Agent) FileStat(fileId string, info *FileInfo) (err error) {
	qs := (protocol.Params{}).With("cid", fileId)
	resp := &webapi.FileStatResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiFileStat, qs, nil, resp); err != nil {
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
func (a *Agent) FileMove(dirId string, fileIds []string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	spec := (&api.FileMoveSpec{}).Init(dirId, fileIds)
	return a.pc.ExecuteApi(spec)
}

// FileCopy copies files into target directory whose id is dirId.
func (a *Agent) FileCopy(dirId string, fileIds []string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	spec := (&api.FileCopySpec{}).Init(dirId, fileIds)
	return a.pc.ExecuteApi(spec)
}

// FileRename renames file to new name.
func (a *Agent) FileRename(fileId, newName string) (err error) {
	spec := (&api.FileRenameSpec{}).Init(fileId, newName)
	return a.pc.ExecuteApi(spec)
}

// FileDelete deletes files.
func (a *Agent) FileDelete(fileIds []string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	spec := (&api.FileDeleteSpec{}).Init(fileIds)
	return a.pc.ExecuteApi(spec)
}

// FileFindDuplications finds all duplicate files which have the same SHA1 hash
// with the given file, and return their ID.
func (a *Agent) FileFindDuplications(fileId string) (dupIds []string, err error) {
	qs := protocol.Params{}.With("file_id", fileId)
	resp := &webapi.BasicResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiFileFindDuplicate, qs, nil, resp); err != nil {
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
