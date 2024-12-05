package protocol

import "github.com/deadblue/elevengo/internal/util"

//lint:ignore U1000 This type is used in generic.
type DirCreateResp struct {
	BasicResp

	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
}

func (r *DirCreateResp) Extract(v *string) (err error) {
	*v = r.FileId
	return
}

//lint:ignore U1000 This type is used in generic.
type DirLocateResp struct {
	BasicResp

	DirId     util.IntNumber `json:"id"`
	IsPrivate util.IntNumber `json:"is_private"`
}

func (r *DirLocateResp) Extract(v *string) (err error) {
	*v = r.DirId.String()
	return
}
