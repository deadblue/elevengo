package api

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/api/base"
	"github.com/deadblue/elevengo/internal/api/errors"
)

type VideoPlayResult struct {
	IsReady       bool
	FileId        string
	FileName      string
	FileSize      int64
	VideoDuration float64
	VideoWidth    int
	VideoHeight   int
	VideoUrl      string
}

//lint:ignore U1000 This type is used in generic.
type _VideoPlayWebResp struct {
	base.BasicResp
	FileId        string           `json:"file_id"`
	ParentId      string           `json:"parent_id"`
	FileName      string           `json:"file_name"`
	FileSize      base.IntNumber   `json:"file_size"`
	FileSha1      string           `json:"sha1"`
	PickCode      string           `json:"pick_code"`
	FileStatus    int              `json:"file_status"`
	VideoDuration base.FloatNumner `json:"play_long"`
	VideoWidth    base.IntNumber   `json:"width"`
	VideoHeight   base.IntNumber   `json:"height"`
	VideoUrl      string           `json:"video_url"`
}

func (r *_VideoPlayWebResp) Extract(v any) error {
	if ptr, ok := v.(*VideoPlayResult); !ok {
		return errors.ErrUnsupportedResult
	} else {
		ptr.IsReady = r.FileStatus == 1
		ptr.FileId = r.FileId
		ptr.FileName = r.FileName
		ptr.FileSize = r.FileSize.Int64()
		ptr.VideoDuration = r.VideoDuration.Float64()
		ptr.VideoWidth = r.VideoWidth.Int()
		ptr.VideoHeight = r.VideoHeight.Int()
		ptr.VideoUrl = r.VideoUrl
	}
	return nil
}

type VideoPlayWebSpec struct {
	base.JsonApiSpec[VideoPlayResult, _VideoPlayWebResp]
}

func (s *VideoPlayWebSpec) Init(pickcode string) *VideoPlayWebSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/video")
	s.QuerySet("pickcode", pickcode)
	return s
}

type _VideoUrl struct {
	Title      string         `json:"title"`
	Definition int            `json:"definition"`
	Width      base.IntNumber `json:"width"`
	Height     base.IntNumber `json:"height"`
	Url        string         `json:"url"`
}

type _VideoPlayPcData struct {
	FileId        string           `json:"file_id"`
	ParentId      string           `json:"parent_id"`
	FileName      string           `json:"file_name"`
	FileSize      base.IntNumber   `json:"file_size"`
	FileSha1      string           `json:"file_sha1"`
	PickCode      string           `json:"pick_code"`
	FileStatus    int              `json:"file_status"`
	VideoDuration base.FloatNumner `json:"play_long"`
	VideoUrls     []*_VideoUrl     `json:"video_url"`
}

type VideoPlayPcSpec struct {
	base.M115ApiSpec[VideoPlayResult]
}

func videoPlayResultExtractor(data []byte, result *VideoPlayResult) (err error) {
	vppd := &_VideoPlayPcData{}
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
	s.M115ApiSpec.Init(
		"https://proapi.115.com/pc/video/play", videoPlayResultExtractor,
	)
	s.ParamSetAll(map[string]string{
		"format":            "app",
		"definition_filter": "1",
		"pickcode":          pickcode,
		"user_id":           userId,
		"appversion":        appVer,
	})
	return s
}
