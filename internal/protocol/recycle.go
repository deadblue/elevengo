package protocol

import (
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/types"
)

//lint:ignore U1000 This type is used in generic.
type RecycleBinListResp struct {
	StandardResp

	Count util.IntNumber `json:"count"`

	Offset int `json:"offset"`
	Limit  int `json:"page_size"`
}

func (r *RecycleBinListResp) Extract(v *types.RecycleBinListResult) (err error) {
	if err = r.StandardResp.Extract(&v.Item); err != nil {
		return
	}
	v.Count = r.Count.Int()
	return
}
