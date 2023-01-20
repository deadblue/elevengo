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
