package elevengo

import (
	"github.com/deadblue/elevengo/internal/core"
	"github.com/deadblue/elevengo/internal/types"
	"strings"
)

const (
	apiDirCreate = "https://webapi.115.com/files/add"
	apiDirGetId  = "https://webapi.115.com/files/getid"
)

// Create a directory under a directory with specific name.
func (a *Agent) DirCreate(parentId, name string) (directoryId string, err error) {
	form := core.NewForm().
		WithString("pid", parentId).
		WithString("cname", name)
	result := &types.FileAddResult{}
	err = a.hc.JsonApi(apiDirCreate, nil, form, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(int(result.ErrorCode), result.Error)
	}
	if err == nil {
		directoryId = result.CategoryId
	}
	return
}

// Retrieve directory Id from full path.
func (a *Agent) DirGetId(path string) (directoryId string, err error) {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	qs := core.NewQueryString().
		WithString("path", path)
	result := &types.DirGetIdResult{}
	err = a.hc.JsonApi(apiDirGetId, qs, nil, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(int(result.ErrorCode), result.Error)
	}
	if err == nil {
		if directoryId = string(result.Id); directoryId == "0" {
			directoryId, err = "", errDirNotExist
		}
	}
	return
}
