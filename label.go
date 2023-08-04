package elevengo

import (
	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/api/errors"
)

type LabelColor int

const (
	LabelNoColor LabelColor = iota
	LabelRed
	LabelOrange
	LabelYellow
	LabelGreen
	LabelBlue
	LabelPurple
	LabelGray

	labelColorMin = LabelNoColor
	labelColorMax = LabelGray
)

type Label struct {
	Id    string
	Name  string
	Color LabelColor
}

// labelIterator implements Iterator[Label].
type labelIterator struct {
	// Offset
	offset int
	// Total count
	count int

	// Cached labels
	labels []*api.LabelInfo
	// Cache index
	index int
	// Cache size
	size int

	// Update function
	uf func(*labelIterator) error
}

func (i *labelIterator) Next() (err error) {
	i.index += 1
	if i.index < i.size {
		return
	}
	i.offset += i.size
	if i.offset >= i.count {
		return errors.ErrReachEnd
	}
	return i.uf(i)
}

func (i *labelIterator) Index() int {
	return i.offset + i.index
}

func (i *labelIterator) Get(label *Label) error {
	if i.index >= i.size {
		return errors.ErrReachEnd
	}
	l := i.labels[i.index]
	label.Id = l.Id
	label.Name = l.Name
	label.Color = LabelColor(api.LabelColorMap[l.Color])
	return nil
}

func (i *labelIterator) Count() int {
	return i.count
}

func (a *Agent) LabelIterate() (it Iterator[Label], err error) {
	li := &labelIterator{
		uf: a.labelIterateInternal,
	}
	if err = a.labelIterateInternal(li); err == nil {
		it = li
	}
	return
}

func (a *Agent) labelIterateInternal(i *labelIterator) (err error) {
	spec := (&api.LabelListSpec{}).Init(i.offset)
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	result := spec.Result
	i.count = result.Total
	i.index, i.size = 0, len(result.List)
	i.labels = make([]*api.LabelInfo, 0, i.size)
	i.labels = append(i.labels, result.List...)
	return
}

// LabelFind finds label whose name is name, and returns it.
// func (a *Agent) LabelFind(name string, label *Label) (err error) {
// 	qs := protocol.Params{}.
// 		With("keyword", name).
// 		WithInt("limit", 10)
// 	resp := &webapi.BasicResponse{}
// 	if err = a.pc.CallJsonApi(webapi.ApiLabelList, qs, nil, resp); err != nil {
// 		return
// 	}
// 	data := &webapi.LabelListData{}
// 	if err = resp.Decode(data); err != nil {
// 		return
// 	}
// 	if data.Total == 0 || data.List[0].Name != name {
// 		err = webapi.ErrNotExist
// 	} else {
// 		label.Id = data.List[0].Id
// 		label.Name = data.List[0].Name
// 		label.Color = LabelColor(webapi.LabelColorMap[data.List[0].Color])
// 	}
// 	return
// }

// LabelCreate creates a label with name and color, returns its ID.
func (a *Agent) LabelCreate(name string, color LabelColor) (labelId string, err error) {
	if color < labelColorMin || color > labelColorMax {
		color = LabelNoColor
	}
	spec := (&api.LabelCreateSpec{}).Init(
		name, api.LabelColors[color],
	)
	if err = a.pc.ExecuteApi(spec); err != nil {
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
	spec := (&api.LabelEditSpec{}).Init(
		label.Id, label.Name,
		api.LabelColors[label.Color],
	)
	return a.pc.ExecuteApi(spec)
}

// LabelDelete deletes a label whose ID is labelId.
func (a *Agent) LabelDelete(labelId string) (err error) {
	if labelId == "" {
		return
	}
	spec := (&api.LabelDeleteSpec{}).Init(labelId)
	return a.pc.ExecuteApi(spec)
}

// FileSetLabels sets labels for a file, you can also remove all labels from it
// by not passing any labelId.
func (a *Agent) FileSetLabels(fileId string, labelIds ...string) (err error) {
	spec := (&api.FileLabelSpec{}).Init(fileId, labelIds)
	return a.pc.ExecuteApi(spec)
}
