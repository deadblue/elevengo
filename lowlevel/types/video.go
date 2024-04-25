package types

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/util"
)

type VideoInfo struct {
	Width   int
	Height  int
	PlayUrl string
}

type VideoPlayResult struct {
	IsReady       bool
	FileId        string
	FileName      string
	FileSize      int64
	VideoDuration float64
	Videos        []*VideoInfo
}

type _VideoUrl struct {
	Title      string         `json:"title"`
	Definition int            `json:"definition"`
	Width      util.IntNumber `json:"width"`
	Height     util.IntNumber `json:"height"`
	Url        string         `json:"url"`
}

type _VideoPlayPcProto struct {
	FileId        string           `json:"file_id"`
	ParentId      string           `json:"parent_id"`
	FileName      string           `json:"file_name"`
	FileSize      util.IntNumber   `json:"file_size"`
	FileSha1      string           `json:"file_sha1"`
	PickCode      string           `json:"pick_code"`
	FileStatus    int              `json:"file_status"`
	VideoDuration util.FloatNumner `json:"play_long"`
	VideoUrls     []*_VideoUrl     `json:"video_url"`
}

func (r *VideoPlayResult) UnmarshalResult(data []byte) (err error) {
	proto := &_VideoPlayPcProto{}
	if err = json.Unmarshal(data, proto); err != nil {
		return
	}
	r.IsReady = proto.FileStatus == 1
	r.FileId = proto.FileId
	r.FileName = proto.FileName
	r.FileSize = proto.FileSize.Int64()
	r.VideoDuration = proto.VideoDuration.Float64()
	r.Videos = make([]*VideoInfo, len(proto.VideoUrls))
	for index, vu := range proto.VideoUrls {
		r.Videos[index] = &VideoInfo{
			Width:   vu.Width.Int(),
			Height:  vu.Height.Int(),
			PlayUrl: vu.Url,
		}
	}
	return nil
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
