package protocol

import (
	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/types"
)

//lint:ignore U1000 This type is used in generic.
type OfflineListResp struct {
	BasicResp

	PageIndex int `json:"page"`
	PageCount int `json:"page_count"`
	PageSize  int `json:"page_row"`

	QuotaTotal  int `json:"total"`
	QuotaRemain int `json:"quota"`

	TaskCount int                  `json:"count"`
	Tasks     []*types.OfflineTask `json:"tasks"`
}

func (r *OfflineListResp) Extract(v any) (err error) {
	if ptr, ok := v.(*types.OfflineListResult); !ok {
		return errors.ErrUnsupportedResult
	} else {
		ptr.PageIndex = r.PageIndex
		ptr.PageCount = r.PageCount
		ptr.PageSize = r.PageSize
		ptr.TaskCount = r.TaskCount
		ptr.Tasks = append(ptr.Tasks, r.Tasks...)
	}
	return nil
}
