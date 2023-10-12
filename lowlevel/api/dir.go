package api

import (
	"github.com/deadblue/elevengo/lowlevel/errors"
)

//lint:ignore U1000 This type is used in generic.
type _DirCreateResp struct {
	_BasicResp

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
	_JsonApiSpec[string, _DirCreateResp]
}

func (s *DirCreateSpec) Init(parentId, name string) *DirCreateSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/files/add")
	s.form.Set("pid", parentId).Set("cname", name)
	return s
}

type DirSetOrderSpec struct {
	_VoidApiSpec
}

func (s *DirSetOrderSpec) Init(dirId string, order string, asc bool) *DirSetOrderSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/order")
	s.form.Set("file_id", dirId).
		Set("fc_mix", "0").
		Set("user_order", order)
	if asc {
		s.form.Set("user_asc", "1")
	} else {
		s.form.Set("user_asc", "0")
	}
	return s
}

//lint:ignore U1000 This type is used in generic.
type _DirLocateResp struct {
	_BasicResp

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
	_JsonApiSpec[string, _DirLocateResp]
}

func (s *DirLocateSpec) Init(path string) *DirLocateSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/files/getid")
	s.query.Set("path", path)
	return s
}
