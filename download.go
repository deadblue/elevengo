package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"net/url"
	"strings"
	"time"
)

const (
	apiFileDownload = "https://webapi.115.com/files/download"
)

// DownloadTicket contains all required information to download a file.
type DownloadTicket struct {
	// Download URL.
	Url string
	// Request headers which SHOULD be sent with download URL.
	Headers map[string]string
	// File name.
	FileName string
	// File size in bytes.
	FileSize int64
}

/*
Create a download ticket.

Now the Agent does not support downloading file, you need using a thirdparty tool to
do that, such as wget/curl/aria2.

TODO: Add curl exmaple.
*/
func (a *Agent) CreateDownloadTicket(pickcode string) (ticket *DownloadTicket, err error) {
	// Get download information
	qs := core.NewQueryString().
		WithString("pickcode", pickcode).
		WithInt64("_", time.Now().Unix())
	result := &internal.DownloadInfoResult{}
	err = a.hc.JsonApi(apiFileDownload, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = internal.MakeFileError(result.MessageCode, result.Message)
	}
	// Create download ticket
	ticket = &DownloadTicket{
		Url:      result.FileUrl,
		Headers:  make(map[string]string),
		FileName: result.FileName,
		FileSize: internal.MustParseInt(result.FileSize),
	}
	// Add user-agent header
	ticket.Headers["User-Agent"] = a.name
	// Add cookie header
	sb := &strings.Builder{}
	downUrl, _ := url.Parse(result.FileUrl)
	for i, ck := range a.cj.Cookies(downUrl) {
		if i > 0 {
			sb.WriteString("; ")
		}
		fmt.Fprintf(sb, "%s=%s", ck.Name, ck.Value)
	}
	ticket.Headers["Cookie"] = sb.String()
	return
}
