package elevengo

import (
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
)

const (
	apiFileVideo = "https://webapi.115.com/files/video"
)

/*
Get HLS data for a video.

For video file, the upstream server support playing online via HLS protocol.
Caller can use this API to get the HLS content, and play it by thirdparty
tools, such as "mpv".

Example:

	// Get HLS content for a video
	hls, err := agent.VideoHlsContent("pickcode")
	if err != nil {
		log.Fatal(err)
	}
	// Start mpv process with reading file from STDIN
	cmd := exec.Command("/path/to/mpv", "-")
	cmd.Stdin = bytes.NewReader(hls)
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
*/
func (a *Agent) VideoHlsContent(pickcode string) (content []byte, err error) {
	// Call video API
	qs := core.NewQueryString().
		WithString("pickcode", pickcode)
	result := &internal.FileVideoResult{}
	err = a.hc.JsonApi(apiFileVideo, qs, nil, result)
	if err == nil {
		if result.IsFailed() {
			err = internal.MakeFileError(result.ErrorCode, result.Error)
		} else if result.FileStatus != 1 {
			err = errVideoNotReady
		}
	}
	if err != nil {
		return
	}
	return a.hc.Get(result.VideoUrl, nil)
}
