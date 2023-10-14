package types

import "github.com/deadblue/elevengo/internal/util"

type LabelInfo struct {
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Color      string         `json:"color"`
	Sort       util.IntNumber `json:"sort"`
	CreateTime int64          `json:"create_time"`
	UpdateTime int64          `json:"update_time"`
}

type LabelListResult struct {
	Total int          `json:"total"`
	List  []*LabelInfo `json:"list"`
	Sort  string       `json:"sort"`
	Order string       `json:"order"`
}

type LabelCreateResult []*LabelInfo
