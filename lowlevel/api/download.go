package api

import (
	"github.com/deadblue/elevengo/lowlevel/types"
)

type DownloadSpec struct {
	_M115ApiSpec[types.DownloadResult]
}

func (s *DownloadSpec) Init(pickcode string) *DownloadSpec {
	s._M115ApiSpec.Init("https://proapi.115.com/app/chrome/downurl", nil)
	s.query.SetNow("t")
	s.params.Set("pickcode", pickcode)
	return s
}
