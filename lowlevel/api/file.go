package api

import (
	"fmt"
	"strings"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
)

const (
	FileOrderByName       = "file_name"
	FileOrderBySize       = "file_size"
	FileOrderByType       = "file_type"
	FileOrderByCreateTime = "user_ptime"
	FileOrderByUpdateTime = "user_utime"
	FileOrderByOpenTime   = "user_otime"
	FileOrderDefault      = FileOrderByCreateTime
)

type FileListSpec struct {
	_JsonApiSpec[types.FileListResult, protocol.FileListResp]

	// Save file order
	fo string
}

func (s *FileListSpec) Init(dirId string, offset, limit int) *FileListSpec {
	s._JsonApiSpec.Init("")
	s.query.Set("format", "json").
		Set("aid", "1").
		Set("cid", dirId).
		Set("show_dir", "1").
		Set("fc_mix", "0").
		SetInt("offset", offset).
		SetInt("limit", limit).
		Set("o", FileOrderDefault).
		Set("asc", "0")
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

type FileSearchSpec struct {
	_JsonApiSpec[types.FileListResult, protocol.FileSearchResp]
}

func (s *FileSearchSpec) Init(offset, limit int) *FileSearchSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/files/search")
	s.query.Set("aid", "1").
		Set("show_dir", "1").
		SetInt("offset", offset).
		SetInt("limit", limit).
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

type FileGetSpec struct {
	_JsonApiSpec[types.FileGetResult, protocol.StandardResp]
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

type FileSetDescSpec struct {
	_VoidApiSpec
}

func (s *FileSetDescSpec) Init(fileId string, desc string) *FileSetDescSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/edit")
	s.form.Set("fid", fileId).
		Set("file_desc", desc)
	return s
}

type FileGetDescSpec struct {
	_JsonApiSpec[string, protocol.FileGetDescResp]
}

func (s *FileGetDescSpec) Init(fileId string) *FileGetDescSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/files/desc")
	s.query.Set("file_id", fileId).
		Set("format", "json").
		Set("compat", "1").
		Set("new_html", "1")
	return s
}

type FileHideSpec struct {
	_VoidApiSpec
}

func (s *FileHideSpec) Init(hide bool, fileIds []string) *FileHideSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/hiddenfiles")
	for i, fileId := range fileIds {
		key := fmt.Sprintf("fid[%d]", i)
		s.form.Set(key, fileId)
	}
	if hide {
		s.form.Set("hidden", "1")
	} else {
		s.form.Set("hidden", "0")
	}
	return s
}
