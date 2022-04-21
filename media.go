package elevengo

import (
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/webapi"
)

const (
	apiFileImage = "https://webapi.115.com/files/image"
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
	qs := (protocol.Params{}).
		With("pickcode", pickcode).
		With("share_id", "0")
	resp := webapi.VideoResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiFileVideo, qs, nil, &resp); err != nil {
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

// ImageUrl gets an image URL which can be embedded into web page.
//func (a *Agent) ImageUrl(pickcode string) (link string, err error) {
//	qs := core.NewQueryString().
//		WithString("pickcode", pickcode).
//		WithInt64("_", time.Now().Unix())
//	result := &types.FileImageResult{}
//	err = a.hc.JsonApi(apiFileImage, qs, nil, result)
//	if err == nil && result.IsFailed() {
//		err = types.MakeFileError(result.ErrorCode, result.Error)
//	}
//	if err == nil {
//		link = result.Data.OriginUrl
//	}
//	return
//}
