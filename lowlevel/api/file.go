package api

import (
	"encoding/json"
	"fmt"
	"strings"

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
	StandardResp

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
	_JsonApiSpec[FileListResult, _FileListResp]

	// Save file order
	fo string
}

func (s *FileListSpec) Init(dirId string, offset int, limit int) *FileListSpec {
	s._JsonApiSpec.Init("")
	s.query.Set("format", "json").
		Set("aid", "1").
		Set("cid", dirId).
		Set("show_dir", "1").
		Set("fc_mix", "0").
		SetInt("offset", offset).
		SetInt("limit", limit)
	// s.QuerySet("snap", "0")
	// s.QuerySet("natsort", "1")
	return s
}

func (s *FileListSpec) Url() string {
	// Select base URL
	if s.fo == FileOrderByName {
		s.baseUrl = "https://aps.115.com/natsort/files.php"
	} else {
		s.baseUrl = "https://webapi.115.com/files"
	}
	return s._JsonApiSpec.Url()
}

func (s *FileListSpec) SetOrder(order string, asc int) {
	s.fo = order
	s.query.Set("o", order).SetInt("asc", asc)
}

func (s *FileListSpec) SetStared() {
	s.query.Set("star", "1")
}

func (s *FileListSpec) SetFileType(fileType int) {
	if fileType != 0 {
		s.query.SetInt("type", fileType)
	}
}

//lint:ignore U1000 This type is used in generic.
type _FileSearchResp struct {
	StandardResp

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
	_JsonApiSpec[FileListResult, _FileSearchResp]
}

func (s *FileSearchSpec) Init(offset int) *FileSearchSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/files/search")
	s.query.Set("aid", "1").
		Set("show_dir", "1").
		SetInt("offset", offset).
		SetInt("limit", FileListLimit).
		Set("format", "json")
	return s
}

func (s *FileSearchSpec) ByKeyword(dirId, keyword string) {
	s.query.Set("cid", dirId).Set("search_value", keyword)
}

func (s *FileSearchSpec) ByLabelId(labelId string) {
	s.query.Set("cid", "0").Set("file_label", labelId)
}

func (s *FileSearchSpec) SetFileType(fileType int) {
	if fileType != 0 {
		s.query.SetInt("type", fileType)
	}
}

type FileGetResult []*FileInfo

type FileGetSpec struct {
	_JsonApiSpec[FileGetResult, StandardResp]
}

func (s *FileGetSpec) Init(fileId string) *FileGetSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/files/get_info")
	s.query.Set("file_id", fileId)
	return s
}

type FileRenameSpec struct {
	_VoidApiSpec
}

func (s *FileRenameSpec) Init() *FileRenameSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/batch_rename")
	return s
}

func (s *FileRenameSpec) Add(fileId, newName string) {
	key := fmt.Sprintf("files_new_name[%s]", fileId)
	s.form.Set(key, newName)
}

type FileMoveSpec struct {
	_VoidApiSpec
}

func (s *FileMoveSpec) Init(dirId string, fileIds []string) *FileMoveSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/move")
	s.form.Set("pid", dirId)
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.form.Set(key, fileId)
	}
	return s
}

type FileCopySpec struct {
	_VoidApiSpec
}

func (s *FileCopySpec) Init(dirId string, fileIds []string) *FileCopySpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/copy")
	s.form.Set("pid", dirId)
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.form.Set(key, fileId)
	}
	return s
}

type FileDeleteSpec struct {
	_VoidApiSpec
}

func (s *FileDeleteSpec) Init(fileIds []string) *FileDeleteSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/rb/delete")
	s.form.Set("ignore_warn", "1")
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.form.Set(key, fileId)
	}
	return s
}

type FileStarSpec struct {
	_VoidApiSpec
}

func (s *FileStarSpec) Init(fileId string, star bool) *FileStarSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/star")
	s.form.Set("file_id", fileId)
	if star {
		s.form.Set("star", "1")
	} else {
		s.form.Set("star", "0")
	}
	return s
}

type FileLabelSpec struct {
	_VoidApiSpec
}

func (s *FileLabelSpec) Init(fileId string, labelIds []string) *FileLabelSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/edit")
	s.form.Set("fid", fileId)
	if len(labelIds) == 0 {
		s.form.Set("file_label", "")
	} else {
		s.form.Set("file_label", strings.Join(labelIds, ","))
	}
	return s
}
