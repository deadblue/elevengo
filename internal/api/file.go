package api

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/api/base"
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
	AreaId     json.Number `json:"aid"`
	CategoryId string      `json:"cid"`
	FileId     string      `json:"fid"`
	ParentId   string      `json:"pid"`

	Name     string      `json:"n"`
	Type     string      `json:"ico"`
	Size     json.Number `json:"s"`
	Sha1     string      `json:"sha"`
	PickCode string      `json:"pc"`

	IsStar json.Number `json:"m"`
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

type _FileListResp struct {
	base.StandardResp

	AreaId     string      `json:"aid"`
	CategoryId json.Number `json:"cid"`
	Count      int         `json:"count"`

	Order  string `json:"order"`
	IsAsc  int    `json:"is_asc"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}

type FileListSpec struct {
	base.JsonApiSpec[_FileListResp, any]
}

func (s *FileListSpec) Init(dirId string) *FileListSpec {
	s.JsonApiSpec.Init("")
	s.QuerySet("format", "json")
	s.QuerySet("aid", "1")
	s.QuerySet("cid", dirId)
	s.QuerySet("show_dir", "1")
	s.QuerySet("fc_mix", "0")
	s.QuerySetInt("offset", 0)
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

func (s *FileListSpec) SetOffset(offset int) {
	s.QuerySetInt("offset", offset)
}

func (s *FileListSpec) SetOrder(order string, asc int) {
	s.QuerySet("o", order)
	s.QuerySetInt("asc", asc)
}
