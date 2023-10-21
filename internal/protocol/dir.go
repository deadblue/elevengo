package protocol

import "github.com/deadblue/elevengo/lowlevel/errors"

//lint:ignore U1000 This type is used in generic.
type DirCreateResp struct {
	BasicResp

	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
}

func (r *DirCreateResp) Extract(v any) (err error) {
	if ptr, ok := v.(*string); !ok {
		err = errors.ErrUnsupportedResult
	} else {
		*ptr = r.FileId
	}
	return
}

//lint:ignore U1000 This type is used in generic.
type DirLocateResp struct {
	BasicResp

	DirId     string `json:"id"`
	IsPrivate string `json:"is_private"`
}

func (r *DirLocateResp) Extract(v any) (err error) {
	if ptr, ok := v.(*string); !ok {
		err = errors.ErrUnsupportedResult
	} else {
		*ptr = r.DirId
	}
	return
}
