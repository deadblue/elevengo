package api

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/api/base"
)

type _VideoPlayWebResp struct {
	base.BasicResp
	FileId        string      `json:"file_id"`
	ParentId      string      `json:"parent_id"`
	FileName      string      `json:"file_name"`
	FileSize      json.Number `json:"file_size"`
	FileSha1      string      `json:"sha1"`
	PickCode      string      `json:"pick_code"`
	VideoStatus   int         `json:"file_status"`
	VideoDuration json.Number `json:"play_long"`
	VideoWidth    json.Number `json:"width"`
	VideoHeight   json.Number `json:"height"`
	VideoUrl      string      `json:"video_url"`
}

type VideoPlayWebSpec struct {
	base.JsonApiSpec[_VideoPlayWebResp]
}

func (s *VideoPlayWebSpec) Init(pickcode string) *VideoPlayWebSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/video")
	s.QuerySet("pickcode", pickcode)
	return s
}

type _VideoInfoPc struct {
	Definition int    `json:"definition"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Url        string `json:"url"`
}

type _VideoPlayPcData struct {
	FileId        string          `json:"file_id"`
	ParentId      string          `json:"parent_id"`
	FileName      string          `json:"file_name"`
	FileSize      json.Number     `json:"file_size"`
	FileSha1      string          `json:"file_sha1"`
	PickCode      string          `json:"pick_code"`
	VideoStatus   int             `json:"file_status"`
	VideoDuration json.Number     `json:"play_long"`
	VideoUrls     []*_VideoInfoPc `json:"video_url"`
}

type VideoPlayPcSpec struct {
	base.M115ApiSpec[_VideoPlayPcData]
}

func (s *VideoPlayPcSpec) Init(userId, appVer, pickcode string) *VideoPlayPcSpec {
	s.M115ApiSpec.Init("https://proapi.115.com/pc/video/play")
	s.ParamSetAll(map[string]string{
		"format":            "app",
		"definition_filter": "1",
		"pickcode":          pickcode,
		"user_id":           userId,
		"appversion":        appVer,
	})
	return s
}
