package elevengo

import (
	"encoding/json"
	"errors"
	"github.com/deadblue/elevengo/internal/crypto/m115"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/gostream/quietly"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	errDownloadNotResult = errors.New("download has no result")
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
DownloadCreateTicket creates ticket which contains all required information to
download a file. Caller can use third-party tools/libraries to download file, such
as wget/curl/aria2.
*/
func (a *Agent) DownloadCreateTicket(pickcode string, ticket *DownloadTicket) (err error) {
	// Generate key for encrypt/decrypt
	key := m115.GenerateKey()

	// Prepare request
	data, _ := json.Marshal(&webapi.DownloadRequest{Pickcode: pickcode})
	qs := protocol.Params{}.WithNow("t")
	form := protocol.Params{}.With("data", m115.Encode(data, key))
	// Send request
	resp := webapi.DownloadResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiDownloadGetUrl, qs, form, &resp); err != nil {
		return
	}
	//
	if !resp.State {
		return errors.New(resp.Message)
	}
	// Parse response
	if data, err = m115.Decode(resp.Data, key); err != nil {
		return
	}
	result := webapi.DownloadResult{}
	if err = json.Unmarshal(data, &result); err != nil {
		return
	}
	if len(result) == 0 {
		return errDownloadNotResult
	}
	for _, v := range result {
		ticket.FileName = v.FileName
		ticket.FileSize, _ = strconv.ParseInt(v.FileSize, 10, 64)
		ticket.Url = v.Url.Url
		ticket.Headers = map[string]string{
			"User-Agent": a.name,
		}
		// Serialize cookie
		cookies := a.pc.ExportCookies(v.Url.Url)
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
		break
	}
	return
}

/*
Download downloads a file from cloud, writes its content into w. If w implements
io.Closer, it will be closed automatically.

This method DOES NOT support multi-thread/resuming, if caller requires those,
use third-party tools/libraries instead.

To monitor the downloading progress, caller can wrap w by
"github.com/deadblue/gostream/observe".
*/
func (a *Agent) Download(pickcode string, w io.Writer) (size int64, err error) {
	if wc, ok := w.(io.WriteCloser); ok {
		defer quietly.Close(wc)
	}

	// Get download ticket.
	ticket := &DownloadTicket{}
	if err = a.DownloadCreateTicket(pickcode, ticket); err != nil {
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
