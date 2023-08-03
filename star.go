package elevengo

import (
	"github.com/deadblue/elevengo/internal/api"
)

// FileStar adds/removes star from a file, whose ID is fileId.
func (a *Agent) FileStar(fileId string, star bool) (err error) {
	spec := (&api.FileStarSpec{}).Init(fileId, star)
	return a.pc.ExecuteApi(spec)
}
