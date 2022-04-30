package elevengo

import (
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"strings"
)

// DirGetId Retrieves directory ID from full path.
func (a *Agent) DirGetId(path string) (dirId string, err error) {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	qs := web.Params{}.With("path", path)
	resp := &webapi.DirGetIdResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiDirGetId, qs, nil, resp); err != nil {
		return
	}
	if err = resp.Err(); err == nil {
		if resp.DirId == "0" {
			err = webapi.ErrNotExist
		} else {
			dirId = string(resp.DirId)
		}
	}
	return
}

type DirOrder int

const (
	DirOrderByTime DirOrder = iota
	DirOrderByType
	DirOrderBySize
	DirOrderByName
)

// DirSetOrder sets how files under it be ordered.
func (a *Agent) DirSetOrder(dirId string, order DirOrder, asc bool) (err error) {
	var orderType string
	switch order {
	case DirOrderByType:
		orderType = "file_type"
	case DirOrderBySize:
		orderType = "file_size"
	case DirOrderByName:
		orderType = "file_name"
	default:
		orderType = "user_ptime"
	}
	form := web.Params{}.
		With("file_id", dirId).
		With("fc_mix", "0").
		With("user_order", orderType)
	if asc {
		form.With("user_asc", "1")
	} else {
		form.With("user_asc", "0")
	}
	resp := &webapi.BasicResponse{}
	return a.wc.CallJsonApi(webapi.ApiDirSetOrder, nil, form, resp)
}
