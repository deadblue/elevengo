package types

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/util"
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

type VideoUrl struct {
	Title      string         `json:"title"`
	Definition int            `json:"definition"`
	Width      util.IntNumber `json:"width"`
	Height     util.IntNumber `json:"height"`
	Url        string         `json:"url"`
}

type VideoPlayPcData struct {
	FileId        string           `json:"file_id"`
	ParentId      string           `json:"parent_id"`
	FileName      string           `json:"file_name"`
	FileSize      util.IntNumber   `json:"file_size"`
	FileSha1      string           `json:"file_sha1"`
	PickCode      string           `json:"pick_code"`
	FileStatus    int              `json:"file_status"`
	VideoDuration util.FloatNumner `json:"play_long"`
	VideoUrls     []*VideoUrl      `json:"video_url"`
}

type _VideoSubtitleProto struct {
	SubtitleId string `json:"sid"`
	Language   string `json:"language"`

	Title string `json:"title"`
	Type  string `json:"type"`
	Url   string `json:"url"`

	SyncTime int `json:"sync_time"`

	IsCaptionMap int    `json:"is_caption_map"`
	CaptionMapId string `json:"caption_map_id"`

	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	PickCode string `json:"pick_code"`
	Sha1     string `json:"sha1"`
}

type VideoSubtitleInfo struct {
	Language string
	Title    string
	Type     string
	Url      string
}

func (i *VideoSubtitleInfo) UnmarshalJSON(data []byte) (err error) {
	if len(data) > 0 && data[0] == '{' {
		proto := &_VideoSubtitleProto{}
		if err = json.Unmarshal(data, proto); err == nil {
			i.Language = proto.Language
			i.Title = proto.Title
			i.Type = proto.Type
			i.Url = proto.Url
		}
	}
	return
}

type VideoSubtitleResult struct {
	AutoLoad VideoSubtitleInfo    `json:"autoload"`
	List     []*VideoSubtitleInfo `json:"list"`
}
