package elevengo

import (
	"os"
	"strings"

	"github.com/deadblue/elevengo/internal/api"
)

// DirMake makes directory under parentId, and returns its ID.
func (a *Agent) DirMake(parentId string, name string) (dirId string, err error) {
	spec := (&api.DirCreateSpec{}).Init(parentId, name)
	if err = a.pc.ExecuteApi(spec); err == nil {
		dirId = spec.Result
	}
	return
}

// DirSetOrder sets how files under the directory be ordered.
func (a *Agent) DirSetOrder(dirId string, order FileOrder, asc bool) (err error) {
	spec := (&api.DirSetOrderSpec{}).Init(
		dirId, getOrderName(order), asc,
	)
	return a.pc.ExecuteApi(spec)
}

// DirGetId retrieves directory ID from full path.
func (a *Agent) DirGetId(path string) (dirId string, err error) {
	path = strings.TrimPrefix(path, "/")
	spec := (&api.DirLocateSpec{}).Init(path)
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	if spec.Result == "0" {
		err = os.ErrNotExist
	} else {
		dirId = spec.Result
	}
	return
}
