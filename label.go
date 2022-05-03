package elevengo

import (
	"fmt"
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
	"strings"
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

func (a *Agent) LabelList() (err error) {
	qs := web.Params{}.
		WithInt("user_id", a.uid).
		With("sort", "create_time").
		With("order", "desc").
		WithInt("offset", 0).
		WithInt("limit", 11500)
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiLabelList, qs, nil, resp); err != nil {
		return
	}
	data := &webapi.LabelListData{}
	if err = resp.Decode(data); err != nil {
		return
	}
	// TODO: How to return?
	return
}

// LabelFind finds label whose name is name, and returns it.
func (a *Agent) LabelFind(name string, label *Label) (err error) {
	qs := web.Params{}.
		With("keyword", name).
		WithInt("limit", 11150)
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiLabelList, qs, nil, resp); err != nil {
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
	form := web.Params{}.
		With("name[]", fmt.Sprintf("%s.%s", name, webapi.LabelColors[color]))
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiLabelAdd, nil, form, resp); err != nil {
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
	form := web.Params{}.
		With("id", label.Id).
		With("name", label.Name).
		With("color", webapi.LabelColors[label.Color])
	return a.wc.CallJsonApi(webapi.ApiLabelEdit, nil, form, &webapi.BasicResponse{})
}

// LabelDelete deletes a label whose ID is labelId.
func (a *Agent) LabelDelete(labelId string) (err error) {
	if labelId == "" {
		return
	}
	form := web.Params{}.With("id", labelId)
	return a.wc.CallJsonApi(webapi.ApiLabelDelete, nil, form, &webapi.BasicResponse{})
}

// FileSetLabel sets labels for a file, you can also remove all labels from it
// by not passing any labelId.
func (a *Agent) FileSetLabel(fileId string, labelIds ...string) (err error) {
	form := web.Params{}.
		With("fid", fileId)
	if len(labelIds) == 0 {
		form.With("file_label", "")
	} else {
		form.With("file_label", strings.Join(labelIds, ","))
	}
	return a.wc.CallJsonApi(webapi.ApiFileEdit, nil, form, &webapi.BasicResponse{})
}
