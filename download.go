package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/gostream/quietly"
	"io"
	"net/http"
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
DownloadCreateTicket creates ticket which contails all necessary information to
download a file. Caller can use thirdparty tools/libraries to perform downloading,
such as wget/curl/aria2.

Example - Downlaod through curl:

	// Create download ticket
	ticket, err := agent.DownloadCreateTicket("pickcode")
	if err != nil {
		log.Fatal(err)
	}
	// Download file via "curl"
	cmd := exec.Command("/usr/bin/curl", ticket.Url)
	for name, value := range ticket.Headers {
		cmd.Args = append(cmd.Args, "-H", fmt.Sprintf("%s: %s", name, value))
	}
	cmd.Args = append(cmd.Args, "-o", ticket.FileName)
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
*/
func (a *Agent) DownloadCreateTicket(pickcode string) (ticket DownloadTicket, err error) {
	// Get download information
	qs := core.NewQueryString().
		WithString("pickcode", pickcode).
		WithInt64("_", time.Now().Unix())
	result := &types.DownloadInfoResult{}
	err = a.hc.JsonApi(apiFileDownload, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(result.MessageCode, result.Message)
	}
	// Create download ticket
	ticket = DownloadTicket{
		Url:      result.FileUrl,
		Headers:  make(map[string]string),
		FileName: result.FileName,
		FileSize: util.MustParseInt(result.FileSize),
	}
	// Add user-agent header
	ticket.Headers["User-Agent"] = a.name
	// Add cookie header
	sb := &strings.Builder{}
	for name, value := range a.hc.Cookies(result.FileUrl) {
		_, _ = fmt.Fprintf(sb, "%s=%s;", name, value)
	}
	ticket.Headers["Cookie"] = sb.String()
	return
}

/*
Download downloads a file from cloud, writes its content into w.
If w implements io.Closer, it will be closed by method.

This method DOSE NOT support multi-thread/resuming, if caller require
those, use thirdparty tools/libraries instead.

To monitor the downloading progress, caller can wrap w by
"github.com/deadblue/gostream/observe".
*/
func (a *Agent) Download(pickcode string, w io.Writer) (size int64, err error) {
	wc, ok := w.(io.WriteCloser)
	if !ok {
		wc = nopWriteCloser{w}
	}
	defer quietly.Close(wc)

	// Get download ticket.
	ticket, err := a.DownloadCreateTicket(pickcode)
	if err != nil {
		return
	}
	// Make download request
	req, err := http.NewRequest(http.MethodGet, ticket.Url, nil)
	if err != nil {
		return
	}
	for name, value := range ticket.Headers {
		req.Header.Set(name, value)
	}
	// Send download request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer quietly.Close(resp.Body)

	// Transfer response body to w
	size, err = io.Copy(w, resp.Body)
	if err == nil && size != ticket.FileSize {
		err = errUnexpectedTransferSize
	}
	return
}

type nopWriteCloser struct {
	io.Writer
}

func (n nopWriteCloser) Close() error {
	return nil
}
