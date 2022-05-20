package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
)

// This files contains deprecated APIs which will be removed in the future.

type FileCursor struct {
	// Transaction
	tx string
	// Cursor parameters
	offset int
	total  int
	// Sort parameters
	order string
	asc   int
}

func (c *FileCursor) checkTransaction(tx string) error {
	if c.tx != tx && c.tx != "" {
		return webapi.ErrInvalidCursor
	} else if c.tx == "" {
		// Initialize cursor
		c.tx = tx
		c.offset = 0
		c.order = "user_ptime"
		c.asc = 0
	}
	return nil
}

func (c *FileCursor) HasMore() bool {
	return c.tx == "" || c.offset < c.total
}

func (c *FileCursor) Total() int {
	return c.total
}

func (c *FileCursor) Remain() int {
	return c.total - c.offset
}

func fileParseListResponse(resp *webapi.FileListResponse, files []*File, cursor *FileCursor) (n int, err error) {
	// Parse response data
	n = len(files)
	result := make([]*webapi.FileInfo, 0, n)
	if err = resp.Decode(&result); err != nil {
		return 0, err
	}
	// Fill files
	if rn := len(result); rn < n {
		n = rn
	}
	for i := 0; i < n; i++ {
		files[i] = (&File{}).from(result[i])
	}
	// Update cursor
	cursor.total = resp.Count
	cursor.offset += n
	return n, nil
}

// FileList lists files list under a directory whose ID is dirId.
// Deprecated: Use FileIterator instead.
func (a *Agent) FileList(dirId string, cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	// Check cursor
	if cursor == nil {
		return 0, webapi.ErrInvalidCursor
	}
	tx := fmt.Sprintf("file_list_%s", dirId)
	if err = cursor.checkTransaction(tx); err != nil {
		return
	}
	// Prepare request
	qs := web.Params{}.
		With("aid", "1").
		With("show_dir", "1").
		With("snap", "0").
		With("natsort", "1").
		With("fc_mix", "0").
		With("format", "json").
		With("cid", dirId).
		With("o", cursor.order).
		WithInt("asc", cursor.asc).
		WithInt("offset", cursor.offset).
		WithInt("limit", n)
	resp := &webapi.FileListResponse{}
	for retry := true; retry; {
		// Select API URL
		apiUrl := webapi.ApiFileList
		if cursor.order == "file_name" {
			apiUrl = webapi.ApiFileListByName
		}
		// Call API
		err = a.wc.CallJsonApi(apiUrl, qs, nil, resp)
		if err == webapi.ErrOrderNotSupport {
			cursor.order, cursor.asc = resp.Order, resp.IsAsc
			qs.With("o", cursor.order).
				WithInt("asc", cursor.asc)
			retry = true
		} else {
			retry = false
		}
	}
	if err != nil {
		return
	}
	// When dirId is not exists, 115 will return the file list under root dir,
	// that should be considered as an error.
	if dirId != string(resp.CategoryId) {
		return 0, webapi.ErrNotExist
	}
	return fileParseListResponse(resp, files, cursor)
}

// FileSearch recursively searches files, whose name contains the keyword and under the directory.
func (a *Agent) FileSearch(dirId, keyword string, cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	// Check cursor
	if cursor == nil {
		return 0, webapi.ErrInvalidCursor
	}
	tx := fmt.Sprintf("file_search_%s_%s", dirId, keyword)
	if err = cursor.checkTransaction(tx); err != nil {
		return
	}
	// Prepare request
	qs := web.Params{}.
		With("aid", "1").
		With("cid", dirId).
		With("search_value", keyword).
		WithInt("offset", cursor.offset).
		WithInt("limit", n).
		With("format", "json")
	resp := &webapi.FileListResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileSearch, qs, nil, resp); err != nil {
		return
	}
	return fileParseListResponse(resp, files, cursor)
}
