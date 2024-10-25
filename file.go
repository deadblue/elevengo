package elevengo

import (
	"context"
	"iter"
	"time"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/types"
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

func (f *File) from(info *types.FileInfo) *File {
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
		f.ModifiedTime = util.ParseFileTime(info.UpdatedTime)
	} else {
		f.ModifiedTime = util.ParseFileTime(info.ModifiedTime)
	}

	f.MediaDuration = info.MediaDuration
	f.IsVideo = info.VideoFlag == 1
	if f.IsVideo {
		f.VideoDefinition = VideoDefinition(info.VideoDefinition)
	}

	return f
}

type fileIterator struct {
	llc client.Client

	// Offset & limit
	offset, limit int
	// Root directory ID
	dirId string
	// Sort order
	order string
	// Is ascending?
	asc int
	// File type
	type_ int
	// File extension
	ext string

	// Iterate mode:
	//  - 1: list
	//  - 2: star
	//  - 3: search
	//  - 4: label
	mode int
	// Parameters for mode=3
	keyword string
	// Parameters for mode=4
	labelId string

	result *types.FileListResult
}

func (i *fileIterator) updateList() (err error) {
	if i.result != nil && i.offset >= i.result.Count {
		return errNoMoreItems
	}
	spec := (&api.FileListSpec{}).Init(i.dirId, i.offset, i.limit)
	if i.mode == 2 {
		spec.SetStared()
	} else {
		if i.order == "" {
			i.order = api.FileOrderDefault
		}
		spec.SetOrder(i.order, i.asc)
	}
	if i.type_ >= 0 {
		spec.SetFileType(i.type_)
	} else {
		spec.SetFileExtension(i.ext)
	}
	for {
		if err = i.llc.CallApi(spec, context.Background()); err == nil {
			break
		}
		if ferr, ok := err.(*errors.FileOrderInvalidError); ok {
			spec.SetOrder(ferr.Order, ferr.Asc)
		} else {
			return
		}
	}
	i.result = &spec.Result
	i.order, i.asc = spec.Result.Order, spec.Result.Asc
	return
}

func (i *fileIterator) updateSearch() (err error) {
	if i.result != nil && i.offset >= i.result.Count {
		return errNoMoreItems
	}
	spec := (&api.FileSearchSpec{}).Init(i.offset, i.limit)
	if i.type_ >= 0 {
		spec.SetFileType(i.type_)
	} else {
		spec.SetFileExtension(i.ext)
	}
	switch i.mode {
	case 3:
		spec.ByKeyword(i.dirId, i.keyword)
	case 4:
		spec.ByLabelId(i.labelId)
	}
	if err = i.llc.CallApi(spec, context.Background()); err == nil {
		i.result = &spec.Result
	}
	return
}

func (i *fileIterator) update() (err error) {
	if i.dirId == "" {
		i.dirId = ""
	}
	if i.limit == 0 {
		i.limit = protocol.FileListLimit
	}
	switch i.mode {
	case 1, 2:
		err = i.updateList()
	case 3, 4:
		err = i.updateSearch()
	}
	return
}

func (i *fileIterator) Count() int {
	if i.result == nil {
		return 0
	}
	return i.result.Count
}

func (i *fileIterator) Items() iter.Seq2[int, *File] {
	return func(yield func(int, *File) bool) {
		for {
			for index, fi := range i.result.Files {
				if cont := yield(i.offset+index, (&File{}).from(fi)); !cont {
					return
				}
			}
			i.offset += i.limit
			if err := i.update(); err != nil {
				break
			}
		}
	}
}

// FileIterate list files under directory, whose id is |dirId|.
func (a *Agent) FileIterate(dirId string) (it Iterator[File], err error) {
	fi := &fileIterator{
		llc:   a.llc,
		dirId: dirId,
		mode:  1,
	}
	if err = fi.update(); err == nil {
		it = fi
	}
	return
}

// FileWithStar lists files with star.
func (a *Agent) FileWithStar(options ...*option.FileListOptions) (it Iterator[File], err error) {
	fi := &fileIterator{
		llc:  a.llc,
		mode: 2,
	}
	// Apply options
	if opts := util.NotNull(options...); opts != nil {
		fi.type_ = opts.Type
		fi.ext = opts.ExtName
	}
	if err = fi.update(); err == nil {
		it = fi
	}
	return
}

// FileSearch recursively searches files under a directory, whose name contains
// the given keyword.
func (a *Agent) FileSearch(
	dirId, keyword string, options ...*option.FileListOptions,
) (it Iterator[File], err error) {
	fi := &fileIterator{
		llc:     a.llc,
		dirId:   dirId,
		mode:    3,
		keyword: keyword,
	}
	// Apply options
	if opts := util.NotNull(options...); opts != nil {
		fi.type_ = opts.Type
		fi.ext = opts.ExtName
	}
	if err = fi.update(); err == nil {
		it = fi
	}
	return
}

// FileLabeled lists files which has specific label.
func (a *Agent) FileWithLabel(
	labelId string, options ...*option.FileListOptions,
) (it Iterator[File], err error) {
	fi := &fileIterator{
		llc:     a.llc,
		mode:    4,
		labelId: labelId,
	}
	// Apply options
	if opts := util.NotNull(options...); opts != nil {
		fi.type_ = opts.Type
		fi.ext = opts.ExtName
	}
	if err = fi.update(); err == nil {
		it = fi
	}
	return
}

// FileGet gets information of a file/directory by its ID.
func (a *Agent) FileGet(fileId string, file *File) (err error) {
	spec := (&api.FileGetSpec{}).Init(fileId)
	if err = a.llc.CallApi(spec, context.Background()); err == nil {
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
	return a.llc.CallApi(spec, context.Background())
}

// FileCopy copies files into target directory whose id is dirId.
func (a *Agent) FileCopy(dirId string, fileIds []string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	spec := (&api.FileCopySpec{}).Init(dirId, fileIds)
	return a.llc.CallApi(spec, context.Background())
}

// FileRename renames file to new name.
func (a *Agent) FileRename(fileId, newName string) (err error) {
	if fileId == "" || newName == "" {
		return
	}
	spec := (&api.FileRenameSpec{}).Init()
	spec.Add(fileId, newName)
	return a.llc.CallApi(spec, context.Background())
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
	return a.llc.CallApi(spec, context.Background())
}

// FileDelete deletes files.
func (a *Agent) FileDelete(fileIds []string) (err error) {
	if len(fileIds) == 0 {
		return
	}
	spec := (&api.FileDeleteSpec{}).Init(fileIds)
	return a.llc.CallApi(spec, context.Background())
}
