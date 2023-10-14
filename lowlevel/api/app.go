package api

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
)

type AppVersionSpec struct {
	_JsonpApiSpec[types.AppVersionResult, protocol.StandardResp]
}

func (s *AppVersionSpec) Init() *AppVersionSpec {
	s._JsonpApiSpec.Init(
		"https://appversion.115.com/1/web/1.0/api/chrome", "get_version",
	)
	return s
}
