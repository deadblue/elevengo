package api

import (
	"github.com/deadblue/elevengo/internal/api/base"
	"github.com/deadblue/elevengo/internal/api/errors"
)

//lint:ignore U1000 This type is used in generic.
type _DirCreateResp struct {
	base.BasicResp
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
}

func (r *_DirCreateResp) Extract(v any) (err error) {
	if ptr, ok := v.(*string); !ok {
		err = errors.ErrUnsupportedResult
	} else {
		*ptr = r.FileId
	}
	return
}

type DirCreateSpec struct {
	base.JsonApiSpec[string, _DirCreateResp]
}

func (s *DirCreateSpec) Init(parentId, name string) *DirCreateSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/add")
	s.FormSet("pid", parentId)
	s.FormSet("cname", name)
	return s
}

type DirSetOrderSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *DirSetOrderSpec) Init(dirId string, order string, asc bool) *DirSetOrderSpec {
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

//lint:ignore U1000 This type is used in generic.
type _DirLocateResp struct {
	base.BasicResp
	DirId     string `json:"id"`
	IsPrivate string `json:"is_private"`
}

func (r *_DirLocateResp) Extract(v any) (err error) {
	if ptr, ok := v.(*string); !ok {
		err = errors.ErrUnsupportedResult
	} else {
		*ptr = r.DirId
	}
	return
}

type DirLocateSpec struct {
	base.JsonApiSpec[string, _DirLocateResp]
}

func (s *DirLocateSpec) Init(path string) *DirLocateSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/getid")
	s.QuerySet("path", path)
	return s
}
