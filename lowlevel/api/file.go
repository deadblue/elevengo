package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/deadblue/elevengo/internal/apibase"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/errors"
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
	AreaId     util.IntNumber `json:"aid"`
	CategoryId string         `json:"cid"`
	FileId     string         `json:"fid"`
	ParentId   string         `json:"pid"`

	Name     string         `json:"n"`
	Type     string         `json:"ico"`
	Size     util.IntNumber `json:"s"`
	Sha1     string         `json:"sha"`
	PickCode string         `json:"pc"`

	IsStar util.Boolean `json:"m"`
	Labels []LabelInfo  `json:"fl"`

	CreatedTime  string `json:"tp"`
	UpdatedTime  string `json:"te"`
	ModifiedTime string `json:"t"`

	// MediaDuration describes duration in seconds for audio/video.
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
	apibase.StandardResp

	AreaId     string         `json:"aid"`
	CategoryId util.IntNumber `json:"cid"`

	Count int `json:"count"`

	Order string `json:"order"`
	IsAsc int    `json:"is_asc"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (r *_FileListResp) Err() (err error) {
	// Handle special error
	if r.ErrorCode2 == errors.CodeFileOrderNotSupported {
		return &errors.ErrFileOrderNotSupported{
			Order: r.Order,
			Asc:   r.IsAsc,
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
	ptr.Order, ptr.Asc = r.Order, r.IsAsc
	return
}

type FileListSpec struct {
	apibase.JsonApiSpec[FileListResult, _FileListResp]
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
	// s.QuerySet("snap", "0")
	// s.QuerySet("natsort", "1")
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

func (s *FileListSpec) SetFileType(fileType int) {
	if fileType != 0 {
		s.QuerySetInt("type", fileType)
	}
}

//lint:ignore U1000 This type is used in generic.
type _FileSearchResp struct {
	apibase.StandardResp

	Folder struct {
		CategoryId string `json:"cid"`
		ParentId   string `json:"pid"`
		Name       string `json:"name"`
	} `json:"folder"`

	Count       int `json:"count"`
	FileCount   int `json:"file_count"`
	FolderCount int `json:"folder_count"`

	Order string `json:"order"`
	IsAsc int    `json:"is_asc"`

	Offset int `json:"offset"`
	Limit  int `json:"page_size"`
}

func (r *_FileSearchResp) Extract(v any) (err error) {
	ptr, ok := v.(*FileListResult)
	if !ok {
		return errors.ErrUnsupportedResult
	}
	ptr.Files = make([]*FileInfo, 0, FileListLimit)
	if err = json.Unmarshal(r.Data, &ptr.Files); err != nil {
		return
	}
	ptr.DirId = r.Folder.CategoryId
	ptr.Count = r.Count
	ptr.Order, ptr.Asc = r.Order, r.IsAsc
	return
}

type FileSearchSpec struct {
	apibase.JsonApiSpec[FileListResult, _FileSearchResp]
}

func (s *FileSearchSpec) Init(offset int) *FileSearchSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/search")
	s.QuerySet("aid", "1")
	s.QuerySet("show_dir", "1")
	s.QuerySetInt("offset", offset)
	s.QuerySetInt("limit", FileListLimit)
	s.QuerySet("format", "json")
	return s
}

func (s *FileSearchSpec) ByKeyword(dirId, keyword string) {
	s.QuerySet("cid", dirId)
	s.QuerySet("search_value", keyword)
}

func (s *FileSearchSpec) ByLabelId(labelId string) {
	s.QuerySet("cid", "0")
	s.QuerySet("file_label", labelId)
}

func (s *FileSearchSpec) SetFileType(fileType int) {
	if fileType != 0 {
		s.QuerySetInt("type", fileType)
	}
}

type FileGetResult []*FileInfo

type FileGetSpec struct {
	apibase.JsonApiSpec[FileGetResult, apibase.StandardResp]
}

func (s *FileGetSpec) Init(fileId string) *FileGetSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/get_info")
	s.QuerySet("file_id", fileId)
	return s
}

type FileRenameSpec struct {
	apibase.JsonApiSpec[apibase.VoidResult, apibase.BasicResp]
}

func (s *FileRenameSpec) Init() *FileRenameSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/batch_rename")
	return s
}

func (s *FileRenameSpec) Add(fileId, newName string) {
	key := fmt.Sprintf("files_new_name[%s]", fileId)
	s.FormSet(key, newName)
}

type FileMoveSpec struct {
	apibase.VoidApiSpec
}

func (s *FileMoveSpec) Init(dirId string, fileIds []string) *FileMoveSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/files/move")
	s.FormSet("pid", dirId)
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.FormSet(key, fileId)
	}
	return s
}

type FileCopySpec struct {
	apibase.VoidApiSpec
}

func (s *FileCopySpec) Init(dirId string, fileIds []string) *FileCopySpec {
	s.VoidApiSpec.Init("https://webapi.115.com/files/copy")
	s.FormSet("pid", dirId)
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.FormSet(key, fileId)
	}
	return s
}

type FileDeleteSpec struct {
	apibase.VoidApiSpec
}

func (s *FileDeleteSpec) Init(fileIds []string) *FileDeleteSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/rb/delete")
	s.FormSet("ignore_warn", "1")
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.FormSet(key, fileId)
	}
	return s
}

type FileStarSpec struct {
	apibase.VoidApiSpec
}

func (s *FileStarSpec) Init(fileId string, star bool) *FileStarSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/files/star")
	s.FormSet("file_id", fileId)
	if star {
		s.FormSet("star", "1")
	} else {
		s.FormSet("star", "0")
	}
	return s
}

type FileLabelSpec struct {
	apibase.VoidApiSpec
}

func (s *FileLabelSpec) Init(fileId string, labelIds []string) *FileLabelSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/files/edit")
	s.FormSet("fid", fileId)
	if len(labelIds) == 0 {
		s.FormSet("file_label", "")
	} else {
		s.FormSet("file_label", strings.Join(labelIds, ","))
	}
	return s
}
