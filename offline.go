package elevengo

import (
	"errors"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"time"
)

const (
	apiOfflineSpace   = "https://115.com/"
	apiOfflineList    = "https://115.com/web/lixian/?ct=lixian&ac=task_lists"
	apiOfflineAddUrl  = "https://115.com/web/lixian/?ct=lixian&ac=add_task_url"
	apiOfflineAddUrls = "https://115.com/web/lixian/?ct=lixian&ac=add_task_urls"
	apiOfflineDelete  = "https://115.com/web/lixian/?ct=lixian&ac=task_del"
	apiOfflineClear   = "https://115.com/web/lixian/?ct=lixian&ac=task_clear"

	errOfflineCaptcha      = 911
	errOfflineTaskExisting = 10008
)

// Parameter for "Client.OfflineClear()" method
type OfflineClearParam struct {
	flag int
}

func (p *OfflineClearParam) All(delete bool) *OfflineClearParam {
	if delete {
		p.flag = 5
	} else {
		p.flag = 1
	}
	return p
}
func (p *OfflineClearParam) Complete(delete bool) *OfflineClearParam {
	if delete {
		p.flag = 4
	} else {
		p.flag = 0
	}
	return p
}
func (p *OfflineClearParam) Failed() *OfflineClearParam {
	p.flag = 2
	return p
}
func (p *OfflineClearParam) Running() *OfflineClearParam {
	p.flag = 3
	return p
}

type OfflineTaskStatus int

func (s OfflineTaskStatus) IsRunning() bool {
	return s == 1
}
func (s OfflineTaskStatus) IsComplete() bool {
	return s == 2
}
func (s OfflineTaskStatus) IsFailed() bool {
	return s == -1
}

type OfflineTask struct {
	InfoHash string
	Name     string
	Url      string
	Status   OfflineTaskStatus
	Percent  int
	FileId   string
}

func (c *Client) offlineSpace() (err error) {
	qs := core.NewQueryString().
		WithString("ct", "offline").
		WithString("ac", "space").
		WithInt64("_", time.Now().Unix())
	result := &internal.OfflineSpaceResult{}
	if err = c.hc.JsonApi(apiOfflineSpace, qs, nil, result); err != nil {
		return
	}
	// store to client
	if c.ot == nil {
		c.ot = &internal.OfflineToken{}
	}
	c.ot.Sign = result.Sign
	c.ot.Time = result.Time
	return nil
}

func (c *Client) callOfflineApi(url string, form core.Form, result interface{}) (err error) {
	if c.ot == nil {
		if err = c.offlineSpace(); err != nil {
			return
		}
	}
	if form == nil {
		form = core.NewForm()
	}
	form.WithString("uid", c.ui.UserId).
		WithString("sign", c.ot.Sign).
		WithInt64("time", c.ot.Time)
	err = c.hc.JsonApi(url, nil, form, result)
	return
}

func (c *Client) OfflineList(page int) (tasks []*OfflineTask, next bool, err error) {
	if page < 1 {
		page = 1
	}
	form := core.NewForm().WithInt("page", page)
	result := &internal.OfflineListResult{}
	err = c.callOfflineApi(apiOfflineList, form, result)
	if err == nil && !result.State {
		err = errors.New(result.ErrorMsg)
	}
	if err != nil {
		return
	}
	tasks = make([]*OfflineTask, len(result.Tasks))
	for index, data := range result.Tasks {
		tasks[index] = &OfflineTask{
			InfoHash: data.InfoHash,
			Name:     data.Name,
			Url:      data.Url,
			Status:   OfflineTaskStatus(data.Status),
			Percent:  data.Precent,
			FileId:   data.FileId,
		}
	}
	next = result.Count-(result.Page*result.PageSize) > 0
	return
}

func (c *Client) OfflineAddUrls(url ...string) (err error) {
	form, isSingle := core.NewForm(), len(url) == 1
	if isSingle {
		form.WithString("url", url[0])
		result := &internal.OfflineAddUrlResult{}
		err = c.callOfflineApi(apiOfflineAddUrl, form, result)
	} else {
		form.WithStrings("url", url)
		result := &internal.OfflineAddUrlsResult{}
		err = c.callOfflineApi(apiOfflineAddUrls, form, result)
	}
	// TODO: return add result
	return
}

func (c *Client) OfflineDelete(deleteFile bool, hash ...string) (err error) {
	form := core.NewForm().WithStrings("hash", hash)
	if deleteFile {
		form.WithInt("flag", 1)
	}
	result := &internal.OfflineBasicResult{}
	err = c.callOfflineApi(apiOfflineDelete, form, result)
	if err == nil && !result.State {
		err = errors.New(result.ErrorMsg)
	}
	return
}

func (c *Client) OfflineClear(params *OfflineClearParam) (err error) {
	if params == nil {
		params = (&OfflineClearParam{}).Complete(false)
	}
	form := core.NewForm().
		WithInt("flag", params.flag)
	result := &internal.OfflineBasicResult{}
	err = c.callOfflineApi(apiOfflineClear, form, result)
	if err == nil && !result.State {
		err = errors.New(result.ErrorMsg)
	}
	return
}
