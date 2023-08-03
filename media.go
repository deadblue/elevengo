package elevengo

import (
	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/api/errors"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
)

// VideoTicket contains all required arguments to play a cloud video.
type VideoTicket struct {
	// Play URL, it is normally a m3u8 URL.
	Url string
	// Request headers which SHOULD be sent with play URL.
	Headers map[string]string
	// File name.
	FileName string
	// File size.
	FileSize int64
	// Video duration in seconds.
	Duration float64
	// Video width.
	Width int
	// Video height.
	Height int
}

// VideoCreateTicket creates a PlayTicket to play the cloud video.
func (a *Agent) VideoCreateTicket(pickcode string, ticket *VideoTicket) (err error) {
	// VideoPlay API for web and PC are different !
	var spec protocol.ApiSpec
	var result *api.VideoPlayResult
	if a.isWeb {
		webSpec := (&api.VideoPlayWebSpec{}).Init(pickcode)
		spec, result = webSpec, &webSpec.Result
	} else {
		pcSpec := (&api.VideoPlayPcSpec{}).Init(
			a.uh.UserId, a.uh.AppVer, pickcode,
		)
		spec, result = pcSpec, &pcSpec.Result
	}
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	if !result.IsReady {
		return errors.ErrVideoNotReady
	}
	ticket.Url = result.VideoUrl
	ticket.Duration = result.VideoDuration
	ticket.Width = result.VideoWidth
	ticket.Height = result.VideoHeight
	ticket.FileName = result.FileName
	ticket.FileSize = result.FileSize
	// Currently(2023-08-02), the play URL for PC does not require any headers,
	// it is extremely recommended to use PC credential.
	if a.isWeb {
		ticket.Headers = map[string]string{
			"User-Agent": a.pc.GetUserAgent(),
			"Cookie":     util.MarshalCookies(a.pc.ExportCookies(ticket.Url)),
		}
	}
	return
}

// ImageGetUrl gets an accessible URL of an image file by its pickcode.
func (a *Agent) ImageGetUrl(pickcode string) (imageUrl string, err error) {
	spec := (&api.ImageGetSpec{}).Init(pickcode)
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	imageUrl = spec.Result.OriginUrl
	return
}
