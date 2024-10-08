package elevengo

import (
	"context"

	"github.com/deadblue/elevengo/lowlevel/api"
)

// FileSetStar adds/removes star from a file, whose ID is fileId.
func (a *Agent) FileSetStar(fileId string, star bool) (err error) {
	spec := (&api.FileStarSpec{}).Init(fileId, star)
	return a.llc.CallApi(spec, context.Background())
}
