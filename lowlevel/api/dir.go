package api

import (
	"github.com/deadblue/elevengo/internal/protocol"
)

type DirCreateSpec struct {
	_JsonApiSpec[string, protocol.DirCreateResp]
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

type DirLocateSpec struct {
	_JsonApiSpec[string, protocol.DirLocateResp]
}

func (s *DirLocateSpec) Init(path string) *DirLocateSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/files/getid")
	s.query.Set("path", path)
	return s
}
