package api

import "github.com/deadblue/elevengo/lowlevel/types"

const (
	LabelListLimit = 30

	LabelColorBlank  = "#000000"
	LabelColorRed    = "#FF4B30"
	LabelColorOrange = "#F78C26"
	LabelColorYellow = "#FFC032"
	LabelColorGreen  = "#43BA80"
	LabelColorBlue   = "#2670FC"
	LabelColorPurple = "#8B69FE"
	LabelColorGray   = "#CCCCCC"
)

type LabelListSpec struct {
	_StandardApiSpec[types.LabelListResult]
}

func (s *LabelListSpec) Init(offset int) *LabelListSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/label/list")
	s.query.SetInt("offset", offset)
	s.query.SetInt("limit", LabelListLimit)
	s.query.Set("sort", "create_time")
	s.query.Set("order", "asc")
	return s
}

type LabelSearchSpec struct {
	_StandardApiSpec[types.LabelListResult]
}

func (s *LabelSearchSpec) Init(keyword string, offset int) *LabelSearchSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/label/list")
	s.query.Set("keyword", keyword)
	s.query.SetInt("offset", offset)
	s.query.SetInt("limit", LabelListLimit)
	return s
}

type LabelCreateSpec struct {
	_StandardApiSpec[types.LabelCreateResult]
}

func (s *LabelCreateSpec) Init(name, color string) *LabelCreateSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/label/add_multi")
	s.form.Set("name[]", name+"\x07"+color)
	return s
}

type LabelEditSpec struct {
	_VoidApiSpec
}

func (s *LabelEditSpec) Init(labelId, name, color string) *LabelEditSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/label/edit")
	s.form.Set("id", labelId).
		Set("name", name).
		Set("color", color)
	return s
}

type LabelDeleteSpec struct {
	_VoidApiSpec
}

func (s *LabelDeleteSpec) Init(labelId string) *LabelDeleteSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/label/delete")
	s.form.Set("id", labelId)
	return s
}

type LabelSetOrderSpec struct {
	_VoidApiSpec
}

func (s *LabelSetOrderSpec) Init(labelId string, order string, asc bool) *LabelSetOrderSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/files/order")
	s.form.Set("module", "label_search").
		Set("file_id", labelId).
		Set("fc_mix", "0").
		Set("user_order", order)
	if asc {
		s.form.Set("user_asc", "1")
	} else {
		s.form.Set("user_asc", "0")
	}
	return s
}
