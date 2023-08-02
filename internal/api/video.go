package api

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/api/base"
	"github.com/deadblue/elevengo/internal/api/errors"
)

type VideoPlayData struct {
	IsReady       bool
	FileId        string
	FileName      string
	FileSize      int64
	VideoDuration float64
	VideoWidth    int
	VideoHeight   int
	VideoUrl      string
}

type _VideoPlayWebResp struct {
	base.BasicResp
	FileId        string      `json:"file_id"`
	ParentId      string      `json:"parent_id"`
	FileName      string      `json:"file_name"`
	FileSize      json.Number `json:"file_size"`
	FileSha1      string      `json:"sha1"`
	PickCode      string      `json:"pick_code"`
	FileStatus    int         `json:"file_status"`
	VideoDuration json.Number `json:"play_long"`
	VideoWidth    json.Number `json:"width"`
	VideoHeight   json.Number `json:"height"`
	VideoUrl      string      `json:"video_url"`
}

func (r *_VideoPlayWebResp) Extract(v any) error {
	if ptr, ok := v.(*VideoPlayData); !ok {
		return errors.ErrUnsupportedData
	} else {
		ptr.IsReady = r.FileStatus == 1
		ptr.FileId = r.FileId
		ptr.FileName = r.FileName
		ptr.FileSize, _ = r.FileSize.Int64()
		ptr.VideoDuration, _ = r.VideoDuration.Float64()
		ptr.VideoWidth = base.MustInt(r.VideoWidth)
		ptr.VideoHeight = base.MustInt(r.VideoHeight)
		ptr.VideoUrl = r.VideoUrl
	}
	return nil
}

type VideoPlayWebSpec struct {
	base.JsonApiSpec[_VideoPlayWebResp, VideoPlayData]
}

func (s *VideoPlayWebSpec) Init(pickcode string) *VideoPlayWebSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/video")
	s.QuerySet("pickcode", pickcode)
	return s
}

type _VideoUrl struct {
	Title      string      `json:"title"`
	Definition int         `json:"definition"`
	Width      json.Number `json:"width"`
	Height     json.Number `json:"height"`
	Url        string      `json:"url"`
}

type _VideoPlayPcData struct {
	FileId        string       `json:"file_id"`
	ParentId      string       `json:"parent_id"`
	FileName      string       `json:"file_name"`
	FileSize      json.Number  `json:"file_size"`
	FileSha1      string       `json:"file_sha1"`
	PickCode      string       `json:"pick_code"`
	FileStatus    int          `json:"file_status"`
	VideoDuration json.Number  `json:"play_long"`
	VideoUrls     []*_VideoUrl `json:"video_url"`
}

type VideoPlayPcSpec struct {
	base.M115ApiSpec[VideoPlayData]
}

func (s *VideoPlayPcSpec) Init(userId, appVer, pickcode string) *VideoPlayPcSpec {
	s.M115ApiSpec.Init("https://proapi.115.com/pc/video/play")
	s.Extractor = videoPlayDataExtractor
	s.ParamSetAll(map[string]string{
		"format":            "app",
		"definition_filter": "1",
		"pickcode":          pickcode,
		"user_id":           userId,
		"appversion":        appVer,
	})
	return s
}

func videoPlayDataExtractor(b []byte, data *VideoPlayData) (err error) {
	vppd := &_VideoPlayPcData{}
	if err = json.Unmarshal(b, vppd); err != nil {
		return
	}
	data.IsReady = vppd.FileStatus == 1
	data.FileId = vppd.FileId
	data.FileName = vppd.FileName
	data.FileSize, _ = vppd.FileSize.Int64()
	data.VideoDuration, _ = vppd.VideoDuration.Float64()
	for _, vu := range vppd.VideoUrls {
		w, h := base.MustInt(vu.Width), base.MustInt(vu.Height)
		if data.VideoWidth < w {
			data.VideoWidth = w
			data.VideoHeight = h
			data.VideoUrl = vu.Url
		}
	}
	return nil
}
