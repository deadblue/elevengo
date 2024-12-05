package protocol

import (
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/types"
)

//lint:ignore U1000 This type is used in generic.
type VideoPlayWebResp struct {
	BasicResp

	FileId        string           `json:"file_id"`
	ParentId      string           `json:"parent_id"`
	FileName      string           `json:"file_name"`
	FileSize      util.IntNumber   `json:"file_size"`
	FileSha1      string           `json:"sha1"`
	PickCode      string           `json:"pick_code"`
	FileStatus    int              `json:"file_status"`
	VideoDuration util.FloatNumner `json:"play_long"`
	VideoWidth    util.IntNumber   `json:"width"`
	VideoHeight   util.IntNumber   `json:"height"`
	VideoUrl      string           `json:"video_url"`
}

func (r *VideoPlayWebResp) Extract(v *types.VideoPlayResult) error {
	v.IsReady = r.FileStatus == 1
	v.FileId = r.FileId
	v.FileName = r.FileName
	v.FileSize = r.FileSize.Int64()
	v.VideoDuration = r.VideoDuration.Float64()
	v.Videos = []*types.VideoInfo{
		{
			Width:   r.VideoWidth.Int(),
			Height:  r.VideoHeight.Int(),
			PlayUrl: r.VideoUrl,
		},
	}
	return nil
}
