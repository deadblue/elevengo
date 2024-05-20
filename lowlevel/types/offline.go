package types

import (
	"encoding/json"

	"github.com/deadblue/elevengo/lowlevel/errors"
)

type OfflineTask struct {
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Url      string `json:"url"`
	AddTime  int64  `json:"add_time"`

	Status     int     `json:"status"`
	Percent    float64 `json:"percentDone"`
	UpdateTime int64   `json:"last_update"`

	FileId string `json:"file_id"`
	DirId  string `json:"wp_path_id"`
}

type OfflineListResult struct {
	PageIndex int
	PageCount int
	PageSize  int

	QuotaTotal  int
	QuotaRemain int

	TaskCount int
	Tasks     []*OfflineTask
}

type _OfflineAddResult struct {
	State   bool   `json:"state"`
	ErrNum  int    `json:"errno"`
	ErrCode int    `json:"errcode"`
	ErrType string `json:"errtype"`

	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

type _OfflineAddUrlsProto struct {
	State   bool `json:"state"`
	ErrNum  int  `json:"errno"`
	ErrCode int  `json:"errcode"`

	Result []*_OfflineAddResult `json:"result"`
}

type OfflineAddUrlsResult []*OfflineTask

func (r *OfflineAddUrlsResult) UnmarshalResult(data []byte) (err error) {
	proto := &_OfflineAddUrlsProto{}
	if err = json.Unmarshal(data, proto); err != nil {
		return
	}
	tasks := make([]*OfflineTask, len(proto.Result))
	for i, r := range proto.Result {
		if r.State || r.ErrCode == errors.CodeOfflineTaskExists {
			tasks[i] = &OfflineTask{}
			tasks[i].InfoHash = r.InfoHash
			tasks[i].Name = r.Name
			tasks[i].Url = r.Url
		} else {
			tasks[i] = nil
		}
	}
	*r = tasks
	return
}
