package elevengo

import (
	"os"
	"strings"

	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
)

type DirOrder int

const (
	DirOrderByName DirOrder = iota
	DirOrderBySize
	DirOrderByType
	DirOrderByCreateTime
	DirOrderByUpdateTime
	DirOrderByOpenTime

	dirOrderMin = DirOrderByName
	dirOrderMax = DirOrderByOpenTime
)

// DirMake makes directory under parentId, and returns its ID.
func (a *Agent) DirMake(parentId string, name string) (dirId string, err error) {
	spec := (&api.DirMakeSpec{}).Init(parentId, name)
	if err = a.pc.ExecuteApi(spec); err == nil {
		dirId = spec.Result
	}
	return
}

// DirSetOrder sets how files under the directory be ordered.
func (a *Agent) DirSetOrder(dirId string, order DirOrder, asc bool) (err error) {
	if order < dirOrderMin || order > dirOrderMax {
		order = DirOrderByUpdateTime
	}
	spec := (&api.DirOrderSpec{}).Init(dirId, api.FileOrderMap[order], asc)
	return a.pc.ExecuteApi(spec)
}

// DirGetId retrieves directory ID from full path.
func (a *Agent) DirGetId(path string) (dirId string, err error) {
	path = strings.TrimPrefix(path, "/")
	qs := protocol.Params{}.With("path", path)
	resp := &webapi.DirLocateResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiDirGetId, qs, nil, resp); err != nil {
		return
	}
	if resp.DirId == "0" {
		err = os.ErrNotExist
	} else {
		dirId = string(resp.DirId)
	}
	return
}
