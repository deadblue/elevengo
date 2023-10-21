package elevengo

import (
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/types"
)

// VideoDefinition values from 115.
type VideoDefinition int

const (
	// Standard Definition, aka. 480P.
	VideoDefinitionSD VideoDefinition = 1
	// High Definition, aka. 720P.
	VideoDefinitionHD VideoDefinition = 2
	// Full-HD, aka. 1080P.
	VideoDefinitionFHD VideoDefinition = 3
	// Another 1080P, what the fuck?
	VideoDefinition1080P VideoDefinition = 4
	// 4K Definition, aka. Ultra-HD.
	VideoDefinition4K VideoDefinition = 5
	// The fallback definition, usually for non-standard resolution.
	VideoDefinitionOrigin VideoDefinition = 100
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
	var spec client.ApiSpec
	var result *types.VideoPlayResult
	if a.isWeb {
		webSpec := (&api.VideoPlayWebSpec{}).Init(pickcode)
		spec, result = webSpec, &webSpec.Result
	} else {
		pcSpec := (&api.VideoPlayPcSpec{}).Init(
			a.uh.UserId, a.uh.AppVer, pickcode,
		)
		spec, result = pcSpec, &pcSpec.Result
	}
	if err = a.llc.CallApi(spec); err != nil {
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
	// Currently(2023-08-02), the play URL for PC does not check any request
	// headers, it is highly recommended to use PC credential.
	if a.isWeb {
		ticket.Headers = map[string]string{
			"User-Agent": a.llc.GetUserAgent(),
			"Cookie":     util.MarshalCookies(a.llc.ExportCookies(ticket.Url)),
		}
	}
	return
}

// ImageGetUrl gets an accessible URL of an image file by its pickcode.
func (a *Agent) ImageGetUrl(pickcode string) (imageUrl string, err error) {
	spec := (&api.ImageGetSpec{}).Init(pickcode)
	if err = a.llc.CallApi(spec); err != nil {
		return
	}
	// The origin URL can be access without cookie.
	imageUrl = spec.Result.OriginUrl
	return
}
