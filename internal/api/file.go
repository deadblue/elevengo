package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deadblue/elevengo/internal/api/base"
	"github.com/deadblue/elevengo/internal/api/errors"
)

const (
	FileOrderByName       = "file_name"
	FileOrderBySize       = "file_size"
	FileOrderByType       = "file_type"
	FileOrderByCreateTime = "user_ptime"
	FileOrderByUpdateTime = "user_utime"
	FileOrderByOpenTime   = "user_otime"
	FileOrderDefault      = FileOrderByCreateTime

	FileListLimit = 32
)

type FileInfo struct {
	AreaId     base.IntNumber `json:"aid"`
	CategoryId string         `json:"cid"`
	FileId     string         `json:"fid"`
	ParentId   string         `json:"pid"`

	Name     string         `json:"n"`
	Type     string         `json:"ico"`
	Size     base.IntNumber `json:"s"`
	Sha1     string         `json:"sha"`
	PickCode string         `json:"pc"`

	IsStar base.Boolean `json:"m"`
	// Labels []*LabelInfo `json:"fl"`

	CreatedTime  string `json:"tp"`
	UpdatedTime  string `json:"te"`
	ModifiedTime string `json:"t"`

	// MediaDuration describes duration in seconds for audio / video.
	MediaDuration float64 `json:"play_long"`

	// Special fields for video
	VideoFlag       int `json:"iv"`
	VideoDefinition int `json:"vdi"`
}

type FileListResult struct {
	DirId  string
	Offset int

	Count int
	Files []*FileInfo

	// Order settings
	Order string
	Asc   int
}

//lint:ignore U1000 This type is used in generic.
type _FileListResp struct {
	base.StandardResp

	AreaId     string         `json:"aid"`
	CategoryId base.IntNumber `json:"cid"`

	Count int `json:"count"`

	Order string `json:"order"`
	Asc   int    `json:"is_asc"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (r *_FileListResp) Err() (err error) {
	// Handle special error
	if r.ErrorCode2 == errors.CodeFileOrderNotSupported {
		return &errors.ErrFileOrderNotSupported{
			Order: r.Order,
			Asc:   r.Asc,
		}
	}
	return r.StandardResp.Err()
}

func (r *_FileListResp) Extract(v any) (err error) {
	ptr, ok := v.(*FileListResult)
	if !ok {
		return errors.ErrUnsupportedResult
	}
	ptr.Files = make([]*FileInfo, 0, FileListLimit)
	if err = json.Unmarshal(r.Data, &ptr.Files); err != nil {
		return
	}
	ptr.DirId = r.CategoryId.String()
	ptr.Count = r.Count
	ptr.Order, ptr.Asc = r.Order, r.Asc
	return
}

type FileListSpec struct {
	base.JsonApiSpec[FileListResult, _FileListResp]
}

func (s *FileListSpec) Init(dirId string, offset int) *FileListSpec {
	s.JsonApiSpec.Init("")
	s.QuerySet("format", "json")
	s.QuerySet("aid", "1")
	s.QuerySet("cid", dirId)
	s.QuerySet("show_dir", "1")
	s.QuerySet("fc_mix", "0")
	s.QuerySetInt("offset", offset)
	s.QuerySetInt("limit", FileListLimit)
	s.QuerySet("o", FileOrderDefault)
	s.QuerySet("asc", "0")
	s.QuerySet("snap", "0")
	s.QuerySet("natsort", "1")
	return s
}

func (s *FileListSpec) Url() string {
	// Select base URL
	order := s.QueryGet("o")
	if order == FileOrderByName {
		s.SetBaseUrl("https://aps.115.com/natsort/files.php")
	} else {
		s.SetBaseUrl("https://webapi.115.com/files")
	}
	return s.JsonApiSpec.Url()
}

func (s *FileListSpec) SetOrder(order string, asc int) {
	s.QuerySet("o", order)
	s.QuerySetInt("asc", asc)
}

func (s *FileListSpec) SetStared() {
	s.QuerySet("star", "1")
}

type FileRenameSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *FileRenameSpec) Init(fileId, newName string) *FileRenameSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/batch_rename")
	if fileId != "" {
		key := fmt.Sprintf("files_new_name[%s]", fileId)
		s.FormSet(key, newName)
	}
	return s
}

type FileMoveSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *FileMoveSpec) Init(dirId string, fileIds []string) *FileMoveSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/move")
	s.FormSet("pid", dirId)
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.FormSet(key, fileId)
	}
	return s
}

type FileCopySpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *FileCopySpec) Init(dirId string, fileIds []string) *FileCopySpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/copy")
	s.FormSet("pid", dirId)
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.FormSet(key, fileId)
	}
	return s
}

type FileDeleteSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *FileDeleteSpec) Init(fileIds []string) *FileDeleteSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/rb/delete")
	s.FormSet("ignore_warn", "1")
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.FormSet(key, fileId)
	}
	return s
}

type FileStarSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *FileStarSpec) Init(fileId string, star bool) *FileStarSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/star")
	s.FormSet("file_id", fileId)
	if star {
		s.FormSet("star", "1")
	} else {
		s.FormSet("star", "0")
	}
	return s
}

type FileLabelSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *FileLabelSpec) Init(fileId string, labelIds []string) *FileLabelSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/edit")
	s.FormSet("fid", fileId)
	if len(labelIds) == 0 {
		s.FormSet("file_label", "")
	} else {
		s.FormSet("file_label", strings.Join(labelIds, ","))
	}
	return s
}
