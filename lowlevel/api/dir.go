package api

import (
	"github.com/deadblue/elevengo/internal/apibase"
	"github.com/deadblue/elevengo/lowlevel/errors"
)

//lint:ignore U1000 This type is used in generic.
type _DirCreateResp struct {
	apibase.BasicResp

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
	apibase.JsonApiSpec[string, _DirCreateResp]
}

func (s *DirCreateSpec) Init(parentId, name string) *DirCreateSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/add")
	s.FormSet("pid", parentId)
	s.FormSet("cname", name)
	return s
}

type DirSetOrderSpec struct {
	apibase.VoidApiSpec
}

func (s *DirSetOrderSpec) Init(dirId string, order string, asc bool) *DirSetOrderSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/files/order")
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
	apibase.BasicResp

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
	apibase.JsonApiSpec[string, _DirLocateResp]
}

func (s *DirLocateSpec) Init(path string) *DirLocateSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/getid")
	s.QuerySet("path", path)
	return s
}
