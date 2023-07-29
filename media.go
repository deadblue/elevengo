package elevengo

import (
	"strings"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
)

// PlayTicket contains all information to play a cloud video.
type PlayTicket struct {
	// Play URL.
	Url string
	// Request headers which SHOULD be used with play URL.
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
func (a *Agent) VideoCreateTicket(pickcode string, ticket *PlayTicket) (err error) {
	qs := protocol.Params{}.
		With("pickcode", pickcode).
		With("share_id", "0").
		With("local", "1")
	resp := &webapi.VideoResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiFileVideo, qs, nil, resp); err != nil {
		return
	}
	if resp.FileStatus != 1 {
		return webapi.ErrVideoNotReady
	}
	ticket.Url = resp.VideoUrl
	ticket.Duration = float64(resp.Duration)
	ticket.Width = int(resp.Width)
	ticket.Height = int(resp.Height)
	ticket.FileName = resp.FileName
	ticket.FileSize = int64(resp.FileSize)
	ticket.Headers = map[string]string{
		"User-Agent": a.pc.GetUserAgent(),
	}
	cookies := a.pc.ExportCookies(ticket.Url)
	if len(cookies) > 0 {
		buf, isFirst := strings.Builder{}, true
		for ck, cv := range cookies {
			if !isFirst {
				buf.WriteString("; ")
			}
			buf.WriteString(ck)
			buf.WriteRune('=')
			buf.WriteString(cv)
			isFirst = false
		}
		ticket.Headers["Cookie"] = buf.String()
	}
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
