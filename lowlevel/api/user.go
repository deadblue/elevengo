package api

import "github.com/deadblue/elevengo/lowlevel/types"

type UserInfoSpec struct {
	_StandardApiSpec[types.UserInfoResult]
}

func (s *UserInfoSpec) Init() *UserInfoSpec {
	s._StandardApiSpec.Init("https://my.115.com/?ct=ajax&ac=nav")
	return s
}
