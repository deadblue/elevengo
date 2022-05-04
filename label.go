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

type labelIterator struct {
	// Offset and total
	o, t int
	// Cached labels
	ls []*webapi.LabelInfo
	// Cache index & cache size
	li, lc int
	// Update function
	updater func(*labelIterator) error
}

func (i *labelIterator) Next() (err error) {
	i.li += 1
	if i.li < i.lc {
		return
	}
	i.o += i.lc
	if i.o >= i.t {
		return webapi.ErrReachEnd
	}
	return i.updater(i)
}

func (i *labelIterator) Get(label *Label) error {
	if i.li >= i.lc {
		return webapi.ErrReachEnd
	}
	l := i.ls[i.li]
	label.Id = l.Id
	label.Name = l.Name
	label.Color = LabelColor(webapi.LabelColorMap[l.Color])
	return nil
}

func (a *Agent) LabelIterate() (it Iterator[Label], err error) {
	li := &labelIterator{
		updater: a.labelIterateInternal,
	}
	if err = a.labelIterateInternal(li); err == nil {
		it = li
	}
	return
}

func (a *Agent) labelIterateInternal(i *labelIterator) (err error) {
	qs := web.Params{}.
		WithInt("user_id", a.uid).
		With("sort", "create_time").
		With("order", "desc").
		WithInt("offset", i.o).
		WithInt("limit", 5)
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiLabelList, qs, nil, resp); err != nil {
		return
	}
	data := &webapi.LabelListData{}
	if err = resp.Decode(data); err == nil {
		i.t = data.Total
		i.li, i.lc = 0, len(data.List)
		// Copy list
		i.ls = make([]*webapi.LabelInfo, 0, i.lc)
		i.ls = append(i.ls, data.List...)
	}
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

// FileSetLabels sets labels for a file, you can also remove all labels from it
// by not passing any labelId.
func (a *Agent) FileSetLabels(fileId string, labelIds ...string) (err error) {
	form := web.Params{}.
		With("fid", fileId)
	if len(labelIds) == 0 {
		form.With("file_label", "")
	} else {
		form.With("file_label", strings.Join(labelIds, ","))
	}
	return a.wc.CallJsonApi(webapi.ApiFileEdit, nil, form, &webapi.BasicResponse{})
}

// FileLabeled lists all files which has specific label.
func (a *Agent) FileLabeled(labelId string, cursor *FileCursor, files []*File) (n int, err error) {
	if n = len(files); n == 0 {
		return
	}
	if cursor == nil {
		return 0, webapi.ErrInvalidCursor
	}
	tx := fmt.Sprintf("file_labeled_%s", labelId)
	if err = cursor.checkTransaction(tx); err != nil {
		return
	}
	// Call API
	qs := web.Params{}.
		With("format", "json").
		With("aid", "1").
		With("cid", "0").
		With("show_dir", "1").
		With("file_label", labelId).
		WithInt("offset", cursor.offset).
		WithInt("limit", n)
	resp := &webapi.FileListResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiFileSearch, qs, nil, resp); err != nil {
		return
	}
	return fileParseListResponse(resp, files, cursor)
}
