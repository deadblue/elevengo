package elevengo

import (
	"github.com/deadblue/elevengo/core"
	"time"
)

type _FileVideoResult struct {
	_FileRetrieveResult
	FileId     string `json:"file_id"`
	FileName   string `json:"file_name"`
	FileSize   string `json:"file_size"`
	FileStatus int    `json:"file_status"`
	Width      string `json:"width"`
	Height     string `json:"height"`
	PlayLong   string `json:"play_long"`
	VideoUrl   string `json:"video_url"`
}

func (c *Client) GetVideoInfo(pickcode string) (hlsInfo []byte, err error) {
	// call API to get video url
	qs := core.NewQueryString().
		WithString("pickcode", pickcode).
		WithInt64("_", time.Now().Unix())
	result := &_FileVideoResult{}
	err = c.requestJson(apiFileVideo, qs, nil, result)
	if err == nil && !result.State {
		err = apiError(result.ErrorNo)
	}
	if err != nil {
		return
	}
	return c.request(result.VideoUrl, nil, nil)
}
