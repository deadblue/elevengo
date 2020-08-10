package elevengo

import (
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"time"
)

const (
	apiFileVideo = "https://webapi.115.com/files/video"
	apiFileImage = "https://webapi.115.com/files/image"
)

/*
Get HLS content of a video.

For video file, the upstream server support HLS streaming. Caller can use this
method to get the HLS content, then play it through thirdparty tools, such as "mpv".
*/
func (a *Agent) VideoHlsContent(pickcode string) (content []byte, err error) {
	// Call video API
	qs := core.NewQueryString().
		WithString("pickcode", pickcode)
	result := &types.FileVideoResult{}
	err = a.hc.JsonApi(apiFileVideo, qs, nil, result)
	if err == nil {
		if result.IsFailed() {
			err = types.MakeFileError(result.ErrorCode, result.Error)
		} else if result.FileStatus != 1 {
			err = errVideoNotReady
		}
	}
	if err != nil {
		return
	}
	return a.hc.Get(result.VideoUrl, nil)
}

// Get a image URL which can be embedded into web page.
func (a *Agent) ImageUrl(pickcode string) (link string, err error) {
	qs := core.NewQueryString().
		WithString("pickcode", pickcode).
		WithInt64("_", time.Now().Unix())
	result := &types.FileImageResult{}
	err = a.hc.JsonApi(apiFileImage, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(result.ErrorCode, result.Error)
	}
	if err == nil {
		link = result.Data.OriginUrl
	}
	return
}
