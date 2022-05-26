package elevengo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deadblue/elevengo/internal/crypto/m115"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"io"
	"strings"
)

var (
	errDownloadNotResult = errors.New("download has no result")

	headerRange = "Range"
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

// DownloadCreateTicket creates ticket which contains all required information
// to download a file. Caller can use third-party tools/libraries to download
// file, such as wget/curl/aria2.
func (a *Agent) DownloadCreateTicket(pickcode string, ticket *DownloadTicket) (err error) {
	// Generate key for encrypt/decrypt
	key := m115.GenerateKey()

	// Prepare request
	data, _ := json.Marshal(&webapi.DownloadRequest{Pickcode: pickcode})
	qs := web.Params{}.WithNow("t")
	form := web.Params{}.With("data", m115.Encode(data, key)).ToForm()
	// Send request
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiDownloadGetUrl, qs, form, resp); err != nil {
		return
	}
	// Parse response
	var resultData string
	if err = resp.Decode(&resultData); err != nil {
		return
	}
	if data, err = m115.Decode(resultData, key); err != nil {
		return
	}
	result := webapi.DownloadData{}
	if err = json.Unmarshal(data, &result); err != nil {
		return
	}
	if len(result) == 0 {
		return errDownloadNotResult
	}
	for _, info := range result {
		a.convertDownloadTicket(info, ticket)
		break
	}
	return
}

func (a *Agent) convertDownloadTicket(info *webapi.DownloadInfo, ticket *DownloadTicket) {
	ticket.FileName = info.FileName
	ticket.FileSize = int64(info.FileSize)
	ticket.Url = info.Url.Url
	ticket.Headers = map[string]string{
		"User-Agent": a.name,
	}
	// Serialize cookie
	cookies := a.wc.ExportCookies(ticket.Url)
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
}

// Get gets content from url using agent underlying HTTP client.
func (a *Agent) Get(url string) (body io.ReadCloser, err error) {
	return a.wc.Get(url, nil, nil)
}

// GetRange gets a part of content from url, which is located by length and offset.
//
// You can use length and offset in 3 cases:
//   - length > 0 and offset < 0: You will get the last length bytes.
//   - length > 0 and offset >= 0: You will get bytes from offset, and at most length bytes.
//   - length < 0 and offset > 0: You will get bytes from offset to the end.
//
// In all other cases, this API equals to Get()
func (a *Agent) GetRange(url string, length, offset int64) (body io.ReadCloser, err error) {
	headers := make(map[string]string)
	// Generate Range header.
	// Reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Range#syntax
	if length > 0 {
		if offset < 0 {
			headers[headerRange] = fmt.Sprintf("bytes=-%d", length)
		} else {
			headers[headerRange] = fmt.Sprintf("bytes=%d-%d", offset, offset+length-1)
		}
	} else if length < 0 {
		if offset > 0 {
			headers[headerRange] = fmt.Sprintf("bytes=%d-", offset)
		}
	}
	return a.wc.Get(url, nil, headers)
}
