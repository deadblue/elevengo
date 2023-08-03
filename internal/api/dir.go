package api

import (
	"github.com/deadblue/elevengo/internal/api/base"
	"github.com/deadblue/elevengo/internal/api/errors"
)

var FileOrderMap = []string{
	FileOrderByName,
	FileOrderBySize,
	FileOrderByType,
	FileOrderByCreateTime,
	FileOrderByUpdateTime,
	FileOrderByOpenTime,
}

//lint:ignore U1000 This type is used in generic.
type _DirMakeResp struct {
	base.BasicResp
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
}

func (r *_DirMakeResp) Extract(v any) (err error) {
	if ptr, ok := v.(*string); !ok {
		err = errors.ErrUnsupportedResult
	} else {
		*ptr = r.FileId
	}
	return
}

type DirMakeSpec struct {
	base.JsonApiSpec[string, _DirMakeResp]
}

func (s *DirMakeSpec) Init(parentId, name string) *DirMakeSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/add")
	s.FormSet("pid", parentId)
	s.FormSet("cname", name)
	return s
}

type DirOrderSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *DirOrderSpec) Init(dirId string, order string, asc bool) *DirOrderSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/order")
	s.FormSetAll(map[string]string{
		"file_id":    dirId,
		"fc_mix":     "0",
		"user_order": order,
	})
	if asc {
		s.FormSet("user_asc", "1")
	} else {
		s.FormSet("user_asc", "0")
	}
	return s
}
