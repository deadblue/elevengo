package api

import "github.com/deadblue/elevengo/lowlevel/types"

type IndexInfoSpec struct {
	_StandardApiSpec[types.IndexInfoResult]
}

func (s *IndexInfoSpec) Init() *IndexInfoSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/files/index_info")
	return s
}
