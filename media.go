package elevengo

import (
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
)

type VideoInfo struct {
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

func (a *Agent) VideoGetInfo(pickcode string, info *VideoInfo) (err error) {
	// Call video API
	qs := web.Params{}.
		With("pickcode", pickcode).
		With("share_id", "0")
	resp := &webapi.VideoResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileVideo, qs, nil, resp); err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	if resp.FileStatus != 1 {
		return webapi.ErrVideoNotReady
	}
	info.FileId = resp.FileId
	info.FileName = resp.FileName
	info.FileSize = int64(resp.FileSize)
	info.FileSha1 = resp.Sha1
	info.PickCode = resp.PickCode
	info.Width = int(resp.Width)
	info.Height = int(resp.Height)
	info.Duration = float64(resp.Duration)
	info.PlayUrl = util.SecretUrl(resp.VideoUrl)
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
	if err = resp.Err(); err != nil {
		return
	}
	// Parse response
	data := &webapi.ImageData{}
	if err = resp.Decode(data); err == nil {
		imageUrl = data.OriginUrl
	}
	return
}
