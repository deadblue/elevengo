package api

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
)

type RecycleBinListSpec struct {
	_JsonApiSpec[types.RecycleBinListResult, protocol.RecycleBinListResp]
}

func (s *RecycleBinListSpec) Init(offset, limit int) *RecycleBinListSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/rb")
	s.query.Set("aid", "7").
		Set("cid", "0").
		Set("format", "json").
		SetInt("offset", offset).
		SetInt("limit", limit)
	return s
}

type RecycleBinCleanSpec struct {
	_VoidApiSpec
}

func (s *RecycleBinCleanSpec) Init(password string) *RecycleBinCleanSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/rb/clean")
	s.form.Set("password", password)
	return s
}
