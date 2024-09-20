package api

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
)

type VideoPlayWebSpec struct {
	_JsonApiSpec[types.VideoPlayResult, protocol.VideoPlayWebResp]
}

func (s *VideoPlayWebSpec) Init(pickcode string) *VideoPlayWebSpec {
	s._JsonApiSpec.Init("https://webapi.115.com/files/video")
	s.query.Set("pickcode", pickcode)
	return s
}

// type VideoPlayPcSpec struct {
// 	_M115ApiSpec[types.VideoPlayResult]
// }

// func (s *VideoPlayPcSpec) Init(pickcode string, common *types.CommonParams) *VideoPlayPcSpec {
// 	s._M115ApiSpec.Init("https://proapi.115.com/pc/video/play")
// 	s.params.Set("format", "app").
// 		Set("appversion", common.AppVer).
// 		Set("user_id", common.UserId).
// 		Set("definition_filter", "1").
// 		Set("pickcode", pickcode)
// 	return s
// }

type VideoSubtitleSpec struct {
	_StandardApiSpec[types.VideoSubtitleResult]
}

func (s *VideoSubtitleSpec) Init(pickcode string) *VideoSubtitleSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/movies/subtitle")
	s.query.Set("pickcode", pickcode)
	return s
}
