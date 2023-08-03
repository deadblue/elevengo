package elevengo

import (
	"fmt"
	"strings"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
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
	// Count
	count int

	// Cached labels
	labels []*webapi.LabelInfo
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
		return webapi.ErrReachEnd
	}
	return i.uf(i)
}

func (i *labelIterator) Index() int {
	return i.offset + i.index
}

func (i *labelIterator) Get(label *Label) error {
	if i.index >= i.size {
		return webapi.ErrReachEnd
	}
	l := i.labels[i.index]
	label.Id = l.Id
	label.Name = l.Name
	label.Color = LabelColor(webapi.LabelColorMap[l.Color])
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
	qs := protocol.Params{}.
		With("user_id", a.uh.UserId).
		With("sort", "create_time").
		With("order", "desc").
		WithInt("offset", i.offset).
		WithInt("limit", 30)
	resp := &webapi.BasicResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiLabelList, qs, nil, resp); err != nil {
		return
	}
	data := &webapi.LabelListData{}
	if err = resp.Decode(data); err == nil {
		i.count = data.Total
		i.index, i.size = 0, len(data.List)
		// Copy list
		i.labels = make([]*webapi.LabelInfo, 0, i.size)
		i.labels = append(i.labels, data.List...)
	}
	return
}

// LabelFind finds label whose name is name, and returns it.
func (a *Agent) LabelFind(name string, label *Label) (err error) {
	qs := protocol.Params{}.
		With("keyword", name).
		WithInt("limit", 10)
	resp := &webapi.BasicResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiLabelList, qs, nil, resp); err != nil {
		return
	}
	data := &webapi.LabelListData{}
	if err = resp.Decode(data); err != nil {
		return
	}
	if data.Total == 0 || data.List[0].Name != name {
		err = webapi.ErrNotExist
	} else {
		label.Id = data.List[0].Id
		label.Name = data.List[0].Name
		label.Color = LabelColor(webapi.LabelColorMap[data.List[0].Color])
	}
	return
}

// LabelCreate creates a label with name and color, returns its ID.
func (a *Agent) LabelCreate(name string, color LabelColor) (labelId string, err error) {
	if color < labelColorMin || color > labelColorMax {
		color = LabelNoColor
	}
	form := protocol.Params{}.
		With("name[]", fmt.Sprintf("%s.%s", name, webapi.LabelColors[color])).
		ToForm()
	resp := &webapi.BasicResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiLabelAdd, nil, form, resp); err != nil {
		return
	}
	var data []*webapi.LabelInfo
	if err = resp.Decode(&data); err == nil {
		if len(data) > 0 {
			labelId = data[0].Id
		} else {
			err = webapi.ErrUnexpected
		}
	}
	return
}

// LabelUpdate updates label's name or color.
func (a *Agent) LabelUpdate(label *Label) (err error) {
	if label == nil || label.Id == "" {
		return
	}
	form := protocol.Params{}.
		With("id", label.Id).
		With("name", label.Name).
		With("color", webapi.LabelColors[label.Color]).
		ToForm()
	return a.pc.CallJsonApi(webapi.ApiLabelEdit, nil, form, &webapi.BasicResponse{})
}

// LabelDelete deletes a label whose ID is labelId.
func (a *Agent) LabelDelete(labelId string) (err error) {
	if labelId == "" {
		return
	}
	form := protocol.Params{}.With("id", labelId).ToForm()
	return a.pc.CallJsonApi(webapi.ApiLabelDelete, nil, form, &webapi.BasicResponse{})
}

// FileSetLabels sets labels for a file, you can also remove all labels from it
// by not passing any labelId.
func (a *Agent) FileSetLabels(fileId string, labelIds ...string) (err error) {
	params := protocol.Params{}.
		With("fid", fileId)
	if len(labelIds) == 0 {
		params.With("file_label", "")
	} else {
		params.With("file_label", strings.Join(labelIds, ","))
	}
	return a.pc.CallJsonApi(webapi.ApiFileEdit, nil, params.ToForm(), &webapi.BasicResponse{})
}
