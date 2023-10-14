package api

import (
	"encoding/json"

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

type VideoPlayPcSpec struct {
	_M115ApiSpec[types.VideoPlayResult]
}

func videoPlayResultExtractor(data []byte, result *types.VideoPlayResult) (err error) {
	vppd := &types.VideoPlayPcData{}
	if err = json.Unmarshal(data, vppd); err != nil {
		return
	}
	result.IsReady = vppd.FileStatus == 1
	result.FileId = vppd.FileId
	result.FileName = vppd.FileName
	result.FileSize = vppd.FileSize.Int64()
	result.VideoDuration = vppd.VideoDuration.Float64()
	for _, vu := range vppd.VideoUrls {
		w, h := vu.Width.Int(), vu.Height.Int()
		if result.VideoWidth < w {
			result.VideoWidth = w
			result.VideoHeight = h
			result.VideoUrl = vu.Url
		}
	}
	return nil
}

func (s *VideoPlayPcSpec) Init(userId, appVer, pickcode string) *VideoPlayPcSpec {
	s._M115ApiSpec.Init(
		"https://proapi.115.com/pc/video/play", videoPlayResultExtractor,
	)
	s.params.Set("format", "app").
		Set("user_id", userId).
		Set("appversion", appVer).
		Set("definition_filter", "1").
		Set("pickcode", pickcode)
	return s
}

type VideoSubtitleSpec struct {
	_StandardApiSpec[types.VideoSubtitleResult]
}

func (s *VideoSubtitleSpec) Init(pickcode string) *VideoSubtitleSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/movies/subtitle")
	s.query.Set("pickcode", pickcode)
	return s
}
