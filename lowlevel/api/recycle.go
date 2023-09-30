package api

import (
	"github.com/deadblue/elevengo/internal/apibase"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/errors"
)

type RecycleBinItem struct {
	FileId     string         `json:"id"`
	FileName   string         `json:"file_name"`
	FileSize   util.IntNumber `json:"file_size"`
	ParentId   string         `json:"cid"`
	ParentName string         `json:"parent_name"`
	DeleteTime util.IntNumber `json:"dtime"`
}

type RecycleBinListResult struct {
	Count int
	Item  []*RecycleBinItem
}

//lint:ignore U1000 This type is used in generic.
type _RecycleBinListResp struct {
	apibase.StandardResp

	Count util.IntNumber `json:"count"`

	Offset int `json:"offset"`
	Limit  int `json:"page_size"`
}

func (r *_RecycleBinListResp) Extract(v any) (err error) {
	ptr, ok := v.(*RecycleBinListResult)
	if !ok {
		return errors.ErrUnsupportedResult
	}
	if err = r.StandardResp.Extract(&ptr.Item); err != nil {
		return
	}
	ptr.Count = r.Count.Int()
	return
}

type RecycleBinListSpec struct {
	apibase.JsonApiSpec[RecycleBinListResult, _RecycleBinListResp]
}

func (s *RecycleBinListSpec) Init(offset int) *RecycleBinListSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/rb")
	s.QuerySet("aid", "7")
	s.QuerySet("cid", "0")
	s.QuerySet("format", "json")
	s.QuerySetInt("offset", offset)
	s.QuerySetInt("limit", FileListLimit)
	return s
}

type RecycleBinCleanSpec struct {
	apibase.VoidApiSpec
}

func (s *RecycleBinCleanSpec) Init(password string) *RecycleBinCleanSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/rb/clean")
	s.FormSet("password", password)
	return s
}
