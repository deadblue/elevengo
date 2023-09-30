package elevengo

import (
	"time"

	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/errors"
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
	Labels []Label

	// Last modified time
	ModifiedTime time.Time

	// Play duration in seconds, for audio/video files.
	MediaDuration float64
	// Is file a video.
	IsVideo bool
	// Definition of the video file.
	VideoDefinition VideoDefinition
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
	f.Labels = make([]Label, len(info.Labels))
	for i, l := range info.Labels {
		f.Labels[i].Id = l.Id
		f.Labels[i].Name = l.Name
		f.Labels[i].Color = labelColorRevMap[l.Color]
	}

	if info.UpdatedTime != "" {
		f.ModifiedTime = api.ParseFileTime(info.UpdatedTime)
	} else {
		f.ModifiedTime = api.ParseFileTime(info.ModifiedTime)
	}

	f.MediaDuration = info.MediaDuration
	f.IsVideo = info.VideoFlag == 1
	if f.IsVideo {
		f.VideoDefinition = VideoDefinition(info.VideoDefinition)
	}

	return f
}

type fileIterator struct {
	// Iterate mode:
	//  - 1: list
	//  - 2: star
	//  - 3: search
	//  - 4: label
	mode int
	// Data offset
	offset int

	// Root dir ID
	dirId string
	// Order field
	order string
	// Order orientation
	asc int
	// Keyword
	keyword string
	// Label Id
	labelId string
	// File type
	fileType int

	// Total count
	count int
	// Cached files
	files []*api.FileInfo
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

// FileIterate list files under directory, whose id is |dirId|.
func (a *Agent) FileIterate(dirId string) (it Iterator[File], err error) {
	fi := &fileIterator{
		mode:   1,
		offset: 0,

		dirId: dirId,
		order: api.FileOrderDefault,
		asc:   0,

		update: a.fileIterateInternal,
	}
	if err = a.fileIterateInternal(fi); err == nil {
		it = fi
	}
	return
}

// FileWithStar lists files with star.
func (a *Agent) FileWithStar(opts ...option.FileListOption) (it Iterator[File], err error) {
	fi := &fileIterator{
		mode:   2,
		offset: 0,

		dirId: "0",

		update: a.fileIterateInternal,
	}
	// Apply options
	for _, opt := range opts {
		switch opt := opt.(type) {
		case option.FileListTypeOption:
			fi.fileType = int(opt)
		}
	}
	if err = a.fileIterateInternal(fi); err == nil {
		it = fi
	}
	return
}

func (a *Agent) fileIterateInternal(fi *fileIterator) (err error) {
	spec := (&api.FileListSpec{}).Init(fi.dirId, fi.offset)
	spec.SetFileType(fi.fileType)
	switch fi.mode {
	case 1:
		spec.SetOrder(fi.order, fi.asc)
	case 2:
		spec.SetStared()
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
	fi.order, fi.asc = spec.Result.Order, spec.Result.Asc
	if fi.count = spec.Result.Count; fi.count > 0 {
		fi.index, fi.size = 0, len(spec.Result.Files)
		fi.files = make([]*api.FileInfo, fi.size)
		copy(fi.files, spec.Result.Files)
	}
	return
}

// FileSearch recursively searches files under a directory, whose name contains
// the given keyword.
func (a *Agent) FileSearch(dirId, keyword string, opts ...option.FileListOption) (it Iterator[File], err error) {
	fi := &fileIterator{
		mode:   3,
		offset: 0,

		dirId:   dirId,
		keyword: keyword,

		update: a.fileSearchInternal,
	}
	// Apply options
	for _, opt := range opts {
		switch opt := opt.(type) {
		case option.FileListTypeOption:
			fi.fileType = int(opt)
		}
	}
	if err = a.fileSearchInternal(fi); err == nil {
		it = fi
	}
	return
}

// FileLabeled lists files which has specific label.
func (a *Agent) FileWithLabel(labelId string, opts ...option.FileListOption) (it Iterator[File], err error) {
	fi := &fileIterator{
		mode:   4,
		offset: 0,

		dirId:   "0",
		labelId: labelId,

		update: a.fileSearchInternal,
	}
	// Apply options
	for _, opt := range opts {
		switch opt := opt.(type) {
		case option.FileListTypeOption:
			fi.fileType = int(opt)
		}
	}
	if err = a.fileSearchInternal(fi); err == nil {
		it = fi
	}
	return
}

func (a *Agent) fileSearchInternal(fi *fileIterator) (err error) {
	spec := (&api.FileSearchSpec{}).Init(fi.offset)
	spec.SetFileType(fi.fileType)
	switch fi.mode {
	case 3:
		spec.ByKeyword(fi.dirId, fi.keyword)
	case 4:
		spec.ByLabelId(fi.labelId)
	}
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	fi.order, fi.asc = spec.Result.Order, spec.Result.Asc
	if fi.count = spec.Result.Count; fi.count > 0 {
		fi.index, fi.size = 0, len(spec.Result.Files)
		fi.files = make([]*api.FileInfo, fi.size)
		copy(fi.files, spec.Result.Files)
	}
	return
}

// FileGet gets information of a file/directory by its ID.
func (a *Agent) FileGet(fileId string, file *File) (err error) {
	spec := (&api.FileGetSpec{}).Init(fileId)
	if err = a.pc.ExecuteApi(spec); err == nil {
		file.from(spec.Result[0])
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
	if fileId == "" || newName == "" {
		return
	}
	spec := (&api.FileRenameSpec{}).Init()
	spec.Add(fileId, newName)
	return a.pc.ExecuteApi(spec)
}

// FileBatchRename renames multiple files.
func (a *Agent) FileBatchRename(nameMap map[string]string) (err error) {
	spec := (&api.FileRenameSpec{}).Init()
	for fileId, newName := range nameMap {
		if fileId == "" || newName == "" {
			continue
		}
		spec.Add(fileId, newName)
	}
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
