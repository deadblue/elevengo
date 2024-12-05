package protocol

import (
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

	TaskCount int               `json:"count"`
	Tasks     []*types.TaskInfo `json:"tasks"`
}

func (r *OfflineListResp) Extract(v *types.OfflineListResult) (err error) {
	v.PageIndex = r.PageIndex
	v.PageCount = r.PageCount
	v.PageSize = r.PageSize
	v.QuotaTotal = r.QuotaTotal
	v.QuotaRemain = r.QuotaRemain
	v.TaskCount = r.TaskCount
	v.Tasks = append(v.Tasks, r.Tasks...)
	return nil
}
