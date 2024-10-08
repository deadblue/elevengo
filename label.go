package elevengo

import (
	"context"
	"iter"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/types"
)

type LabelColor int

const (
	LabelColorBlank LabelColor = iota
	LabelColorRed
	LabelColorOrange
	LabelColorYellow
	LabelColorGreen
	LabelColorBlue
	LabelColorPurple
	LabelColorGray
)

var (
	labelColorMap = map[LabelColor]string{
		LabelColorBlank:  api.LabelColorBlank,
		LabelColorRed:    api.LabelColorRed,
		LabelColorOrange: api.LabelColorOrange,
		LabelColorYellow: api.LabelColorYellow,
		LabelColorGreen:  api.LabelColorGreen,
		LabelColorBlue:   api.LabelColorBlue,
		LabelColorPurple: api.LabelColorPurple,
		LabelColorGray:   api.LabelColorGray,
	}

	labelColorRevMap = map[string]LabelColor{
		api.LabelColorBlank:  LabelColorBlank,
		api.LabelColorRed:    LabelColorRed,
		api.LabelColorOrange: LabelColorOrange,
		api.LabelColorYellow: LabelColorYellow,
		api.LabelColorGreen:  LabelColorGreen,
		api.LabelColorBlue:   LabelColorBlue,
		api.LabelColorPurple: LabelColorPurple,
		api.LabelColorGray:   LabelColorGray,
	}
)

type Label struct {
	Id    string
	Name  string
	Color LabelColor
}

func (l *Label) from(info *types.LabelInfo) *Label {
	l.Id = info.Id
	l.Name = info.Name
	l.Color = labelColorRevMap[info.Color]
	return l
}

type labelIterator struct {
	llc    client.Client
	offset int
	limit  int
	result *types.LabelListResult
}

func (i *labelIterator) update() (err error) {
	if i.result != nil && i.offset >= i.result.Total {
		return errNoMoreItems
	}
	spec := (&api.LabelListSpec{}).Init(i.offset, i.limit)
	if err = i.llc.CallApi(spec, context.Background()); err == nil {
		i.result = &spec.Result
	}
	return
}

func (i *labelIterator) Count() int {
	if i.result == nil {
		return 0
	}
	return i.result.Total
}

func (i *labelIterator) Items() iter.Seq2[int, *Label] {
	return func(yield func(int, *Label) bool) {
		for {
			for index, li := range i.result.List {
				if stop := !yield(i.offset+index, (&Label{}).from(li)); stop {
					return
				}
			}
			i.offset += i.limit
			if err := i.update(); err != nil {
				return
			}
		}
	}
}

func (a *Agent) LabelIterate() (it Iterator[*Label], err error) {
	li := &labelIterator{
		llc:    a.llc,
		offset: 0,
		limit:  protocol.LabelListLimit,
	}
	if err = li.update(); err == nil {
		it = li
	}
	return
}

// LabelFind finds label whose name is name, and returns it.
func (a *Agent) LabelFind(name string, label *Label) (err error) {
	spec := (&api.LabelSearchSpec{}).Init(name, 0, protocol.LabelListLimit)
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}

	if spec.Result.Total == 0 || spec.Result.List[0].Name != name {
		err = errors.ErrNotExist
	} else {
		li := spec.Result.List[0]
		label.Id = li.Id
		label.Name = li.Name
		label.Color = labelColorRevMap[li.Color]
	}
	return
}

// LabelCreate creates a label with name and color, returns its ID.
func (a *Agent) LabelCreate(name string, color LabelColor) (labelId string, err error) {
	colorName, ok := labelColorMap[color]
	if !ok {
		colorName = api.LabelColorBlank
	}
	spec := (&api.LabelCreateSpec{}).Init(
		name, colorName,
	)
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}
	if len(spec.Result) > 0 {
		labelId = spec.Result[0].Id
	}
	return
}

// LabelUpdate updates label's name or color.
func (a *Agent) LabelUpdate(label *Label) (err error) {
	if label == nil || label.Id == "" {
		return
	}
	colorName, ok := labelColorMap[label.Color]
	if !ok {
		colorName = api.LabelColorBlank
	}
	spec := (&api.LabelEditSpec{}).Init(
		label.Id, label.Name, colorName,
	)
	return a.llc.CallApi(spec, context.Background())
}

// LabelDelete deletes a label whose ID is labelId.
func (a *Agent) LabelDelete(labelId string) (err error) {
	if labelId == "" {
		return
	}
	spec := (&api.LabelDeleteSpec{}).Init(labelId)
	return a.llc.CallApi(spec, context.Background())
}

func (a *Agent) LabelSetOrder(labelId string, order FileOrder, asc bool) (err error) {
	spec := (&api.LabelSetOrderSpec{}).Init(
		labelId, getOrderName(order), asc,
	)
	return a.llc.CallApi(spec, context.Background())
}

// FileSetLabels sets labels for a file, you can also remove all labels from it
// by not passing any labelId.
func (a *Agent) FileSetLabels(fileId string, labelIds ...string) (err error) {
	spec := (&api.FileLabelSpec{}).Init(fileId, labelIds)
	return a.llc.CallApi(spec, context.Background())
}
