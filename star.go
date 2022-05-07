package elevengo

import (
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
)

// FileStar adds/removes star from a file, whose ID is fileId.
func (a *Agent) FileStar(fileId string, star bool) (err error) {
	form := web.Params{}.
		With("file_id", fileId).
		WithInt("star", webapi.BoolToInt(star)).
		ToForm()
	return a.wc.CallJsonApi(webapi.ApiFileStar, nil, form, &webapi.BasicResponse{})
}

// FileStared lists all stared files.
func (a *Agent) FileStared(cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	// Check cursor
	if cursor == nil {
		return 0, webapi.ErrInvalidCursor
	}
	tx := "file_stared"
	if err = cursor.checkTransaction(tx); err != nil {
		return
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
	return fileParseListResponse(resp, files, cursor)
}
