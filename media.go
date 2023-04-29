package elevengo

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/webapi"
)

// Video contains information of a video file on cloud.
type Video struct {

	// File ID
	FileId string
	// File Name
	FileName string
	// File size in bytes
	FileSize int64
	// Pick code for downloading
	PickCode string
	// SHA-1 hash
	Sha1 string

	// Video width
	Width int
	// Video height
	Height int
	// Video duration in seconds
	Duration float64

	// Play URL, usually is a m3u8 URL
	PlayUrl string
}

// VideoGet gets information of a video file by its pickcode.
func (a *Agent) VideoGet(pickcode string, video *Video) (err error) {
	// Call video API
	qs := protocol.Params{}.
		With("pickcode", pickcode).
		With("share_id", "0")
	resp := &webapi.VideoResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiFileVideo, qs, nil, resp); err != nil {
		return
	}
	if resp.FileStatus != 1 {
		return webapi.ErrVideoNotReady
	}
	video.FileId = resp.FileId
	video.FileName = resp.FileName
	video.FileSize = int64(resp.FileSize)
	video.PickCode = resp.PickCode
	video.Sha1 = resp.Sha1
	video.Width = int(resp.Width)
	video.Height = int(resp.Height)
	video.Duration = float64(resp.Duration)
	video.PlayUrl = util.SecretUrl(resp.VideoUrl)
	return
}

// ImageGetUrl gets an accessible URL of an image file by its pickcode.
func (a *Agent) ImageGetUrl(pickcode string) (imageUrl string, err error) {
	qs := protocol.Params{}.
		With("pickcode", pickcode).
		WithNow("_")
	resp := &webapi.BasicResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiFileImage, qs, nil, resp); err != nil {
		return
	}
	// Parse response
	data := &webapi.ImageData{}
	if err = resp.Decode(data); err == nil {
		imageUrl = data.OriginUrl
	}
	return
}
