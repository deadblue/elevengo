package elevengo

import (
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
)

// FileStar adds/removes star from a file, whose ID is fileId.
func (a *Agent) FileStar(fileId string, star bool) (err error) {
	form := web.Params{}.
		With("file_id", fileId).
		WithInt("star", webapi.BoolToInt(star))
	return a.wc.CallJsonApi(webapi.ApiFileStar, nil, form, &webapi.BasicResponse{})
}

// FileListStared lists all started files.
func (a *Agent) FileListStared(cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	// Initialize cursor
	if !cursor.init {
		cursor.init = true
	}
	qs := web.Params{}.
		With("aid", "1").
		With("cid", "0").
		With("show_dir", "1").
		With("star", "1").
		With("format", "json").
		WithInt("offset", cursor.offset).
		WithInt("limit", n)
	resp := &webapi.FileListResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileList, qs, nil, resp); err != nil {
		return
	}
	// Parse result
	result := make([]*webapi.FileInfo, 0, n)
	if err = resp.Decode(&result); err != nil {
		return
	}
	// Fill result to files
	if rn := len(result); rn < n {
		n = rn
	}
	for i := 0; i < n; i++ {
		files[i] = (&File{}).from(result[i])
	}
	// Update cursor
	cursor.offset += n
	cursor.total = resp.Count
	return
}
