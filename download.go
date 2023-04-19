package elevengo

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/deadblue/elevengo/internal/crypto/m115"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
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
	if !result.IsValid() {
		return webapi.ErrDownloadEmpty
	}
	for _, info := range result {
		if info.FileSize == 0 {
			err = webapi.ErrDownloadDirectory
		} else {
			a.convertDownloadTicket(info, ticket)
		}
		break
	}
	return
}

func (a *Agent) convertDownloadTicket(info *webapi.DownloadInfo, ticket *DownloadTicket) {
	ticket.FileName = info.FileName
	ticket.FileSize = int64(info.FileSize)
	ticket.Url = info.Url.Url
	ticket.Headers = map[string]string{
		"User-Agent": a.wc.GetUserAgent(),
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

// Range is used in Agent.GetRange().
type Range struct {
	start, end int64
}

func (r *Range) headerValue() string {
	// Generate Range header.
	// Reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Range#syntax
	if r.start < 0 {
		return fmt.Sprintf("bytes=%d", r.start)
	} else {
		if r.end < 0 {
			return fmt.Sprintf("bytes=%d-", r.start)
		} else if r.end > r.start {
			return fmt.Sprintf("bytes=%d-%d", r.start, r.end)
		}
	}
	// (r.start >= 0 && r.end <= r.start) is an invalid range
	return ""
}

// RangeFirst makes a Range parameter to request the first `length` bytes.
func RangeFirst(length int64) Range {
	return Range{
		start: 0,
		end:   length - 1,
	}
}

// RangeLast makes a Range parameter to request the last `length` bytes.
func RangeLast(length int64) Range {
	return Range{
		start: 0 - length,
		end:   0,
	}
}

// RangeMiddle makes a Range parameter to request content starts from `offset`,
// and has `length` bytes (at most).
//
// You can pass a negative number in `length`, to request content starts from
// `offset` to the end.
func RangeMiddle(offset, length int64) Range {
	end := offset + length - 1
	if length < 0 {
		end = -1
	}
	return Range{
		start: offset,
		end:   end,
	}
}

// GetRange gets partial content from |url|, which is located by |rng|.
func (a *Agent) GetRange(url string, rng Range) (body io.ReadCloser, err error) {
	headers := make(map[string]string)
	if value := rng.headerValue(); value != "" {
		headers["Range"] = value
	}
	return a.wc.Get(url, nil, headers)
}
