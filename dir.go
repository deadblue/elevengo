package elevengo

import (
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"os"
	"strings"
)

type DirOrder int

const (
	DirOrderByTime DirOrder = iota
	DirOrderByType
	DirOrderBySize
	DirOrderByName

	dirOrderMin = DirOrderByTime
	dirOrderMax = DirOrderByName
)

// DirMake makes directory under parentId, and returns its ID.
func (a *Agent) DirMake(parentId string, name string) (dirId string, err error) {
	form := web.Params{}.
		With("pid", parentId).
		With("cname", name).ToForm()
	resp := &webapi.DirMakeResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiDirAdd, nil, form, resp); err == nil {
		dirId = resp.CategoryId
	}
	return
}

// DirSetOrder sets how files under the directory be ordered.
func (a *Agent) DirSetOrder(dirId string, order DirOrder, asc bool) (err error) {
	if order < dirOrderMin || order > dirOrderMax {
		order = DirOrderByTime
	}
	form := web.Params{}.
		With("file_id", dirId).
		With("fc_mix", "0").
		With("user_order", webapi.DirOrderModes[order]).
		WithInt("user_asc", webapi.BoolToInt(asc)).
		ToForm()
	return a.wc.CallJsonApi(webapi.ApiDirSetOrder, nil, form, &webapi.BasicResponse{})
}

// DirGetId retrieves directory ID from full path.
func (a *Agent) DirGetId(path string) (dirId string, err error) {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	qs := web.Params{}.With("path", path)
	resp := &webapi.DirLocateResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiDirGetId, qs, nil, resp); err != nil {
		return
	}
	if resp.DirId == "0" {
		err = os.ErrNotExist
	} else {
		dirId = string(resp.DirId)
	}
	return
}
