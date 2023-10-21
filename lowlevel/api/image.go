package api

import "github.com/deadblue/elevengo/lowlevel/types"

type ImageGetSpec struct {
	_StandardApiSpec[types.ImageGetResult]
}

func (s *ImageGetSpec) Init(pickcode string) *ImageGetSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/files/image")
	s.query.Set("pickcode", pickcode)
	s.query.SetNow("_")
	return s
}
