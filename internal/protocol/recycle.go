package protocol

import (
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/types"
)

//lint:ignore U1000 This type is used in generic.
type RecycleBinListResp struct {
	StandardResp

	Count util.IntNumber `json:"count"`

	Offset int `json:"offset"`
	Limit  int `json:"page_size"`
}

func (r *RecycleBinListResp) Extract(v any) (err error) {
	ptr, ok := v.(*types.RecycleBinListResult)
	if !ok {
		return errors.ErrUnsupportedResult
	}
	if err = r.StandardResp.Extract(&ptr.Item); err != nil {
		return
	}
	ptr.Count = r.Count.Int()
	return
}
