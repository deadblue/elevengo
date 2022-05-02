package elevengo

import (
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
)

type Video struct {
	FileId   string
	FileName string
	FileSize int64
	FileSha1 string
	PickCode string
	Width    int
	Height   int
	Duration float64
	PlayUrl  string
}

func (a *Agent) VideoGetInfo(pickcode string, video *Video) (err error) {
	// Call video API
	qs := web.Params{}.
		With("pickcode", pickcode).
		With("share_id", "0")
	resp := &webapi.VideoResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileVideo, qs, nil, resp); err != nil {
		return
	}
	if resp.FileStatus != 1 {
		return webapi.ErrVideoNotReady
	}
	video.FileId = resp.FileId
	video.FileName = resp.FileName
	video.FileSize = int64(resp.FileSize)
	video.FileSha1 = resp.Sha1
	video.PickCode = resp.PickCode
	video.Width = int(resp.Width)
	video.Height = int(resp.Height)
	video.Duration = float64(resp.Duration)
	video.PlayUrl = util.SecretUrl(resp.VideoUrl)
	return
}

// ImageGetUrl gets an accessible image URL of given pickcode, which is from an image file.
func (a *Agent) ImageGetUrl(pickcode string) (imageUrl string, err error) {
	qs := web.Params{}.
		With("pickcode", pickcode).
		WithNow("_")
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileImage, qs, nil, resp); err != nil {
		return
	}
	// Parse response
	data := &webapi.ImageData{}
	if err = resp.Decode(data); err == nil {
		imageUrl = data.OriginUrl
	}
	return
}
