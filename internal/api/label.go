package api

import "github.com/deadblue/elevengo/internal/api/base"

const (
	LabelListLimit = 30
)

var (
	LabelColors = []string{
		// No Color
		"#000000",
		// Red
		"#FF4B30",
		// Orange
		"#F78C26",
		// Yellow
		"#FFC032",
		// Green
		"#43BA80",
		// Blue
		"#2670FC",
		// Purple
		"#8B69FE",
		// Gray
		"#CCCCCC",
	}

	LabelColorMap = map[string]int{
		"#000000": 0,
		"#FF4B30": 1,
		"#F78C26": 2,
		"#FFC032": 3,
		"#43BA80": 4,
		"#2670FC": 5,
		"#8B69FE": 6,
		"#CCCCCC": 7,
	}
)

type LabelInfo struct {
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Color      string         `json:"color"`
	Sort       base.IntNumber `json:"sort"`
	CreateTime int64          `json:"create_time"`
	UpdateTime int64          `json:"update_time"`
}

type LabelListResult struct {
	Total int          `json:"total"`
	List  []*LabelInfo `json:"list"`
	Sort  string       `json:"sort"`
	Order string       `json:"order"`
}

type LabelListSpec struct {
	base.JsonApiSpec[LabelListResult, base.StandardResp]
}

func (s *LabelListSpec) Init(offset int) *LabelListSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/label/list")
	s.QuerySetInt("offset", offset)
	s.QuerySetInt("limit", LabelListLimit)
	s.QuerySet("sort", "create_time")
	s.QuerySet("order", "asc")
	return s
}

type LabelCreateResult []*LabelInfo

type LabelCreateSpec struct {
	base.JsonApiSpec[LabelCreateResult, base.StandardResp]
}

func (s *LabelCreateSpec) Init(name, color string) *LabelCreateSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/label/add_multi")
	s.FormSet("name[]", name+"\x07"+color)
	return s
}

type LabelEditSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *LabelEditSpec) Init(labelId, name, color string) *LabelEditSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/label/edit")
	s.FormSetAll(map[string]string{
		"id":    labelId,
		"name":  name,
		"color": color,
	})
	return s
}

type LabelDeleteSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *LabelDeleteSpec) Init(labelId string) *LabelDeleteSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/label/delete")
	s.FormSet("id", labelId)
	return s
}
