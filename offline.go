package elevengo

import (
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
)

type offlineCursor struct {
	used      bool
	page      int
	pageCount int
	total     int
}

func (c *offlineCursor) HasMore() bool {
	return !c.used || c.page < c.pageCount
}
func (c *offlineCursor) Next() {
	c.page += 1
}
func (c *offlineCursor) Total() int {
	return c.total
}

// Create a cursor for "Agent.OfflineList()" method.
func OfflineCursor() Cursor {
	return &offlineCursor{
		used:      false,
		page:      1,
		pageCount: 0,
		total:     0,
	}
}

// Parameter for "Agent.OfflineClear()" method.
// Default value is to clear all done tasks without deleteing downloaded files.
type OfflineClearFlag struct {
	flag int
}

// Clear all tasks, delete the downloaded files if "delete" is true.
func (f *OfflineClearFlag) All(delete bool) *OfflineClearFlag {
	if delete {
		f.flag = 5
	} else {
		f.flag = 1
	}
	return f
}

// Clear done tasks, delete the downloaded files if "delete" is true.
func (f *OfflineClearFlag) Done(delete bool) *OfflineClearFlag {
	if delete {
		f.flag = 4
	} else {
		f.flag = 0
	}
	return f
}

// Clear failed tasks.
func (f *OfflineClearFlag) Failed() *OfflineClearFlag {
	f.flag = 2
	return f
}

// Clean running tasks.
func (f *OfflineClearFlag) Running() *OfflineClearFlag {
	f.flag = 3
	return f
}

// Describe status of an offline task.
type OfflineTaskStatus int

// Return true if the task is still running.
func (s OfflineTaskStatus) IsRunning() bool {
	return s == 1
}

// Return true if the task has been done.
func (s OfflineTaskStatus) IsDone() bool {
	return s == 2
}

// Return true if the task has been failed.
func (s OfflineTaskStatus) IsFailed() bool {
	return s == -1
}

// OfflineTask describe a remote download task.
type OfflineTask struct {
	// Unique hash of the task.
	InfoHash string
	// Task name.
	Name string
	// Task URL.
	Url string
	// Task status.
	Status OfflineTaskStatus
	// Download percent of the task, 0 to 100.
	Percent int
	// File ID of the downloaded file on remote server.
	FileId string
}

func (a *Agent) updateOfflineToken() (err error) {
	qs := core.NewQueryString().
		WithString("ct", "offline").
		WithString("ac", "space").
		WithInt64("_", time.Now().Unix())
	result := &internal.OfflineSpaceResult{}
	if err = a.hc.JsonApi(apiOfflineSpace, qs, nil, result); err != nil {
		return
	}
	// store to client
	if a.ot == nil {
		a.ot = &internal.OfflineToken{}
	}
	a.ot.Sign = result.Sign
	a.ot.Time = result.Time
	return nil
}

func (a *Agent) callOfflineApi(url string, form core.Form, result interface{}) (err error) {
	if a.ot == nil {
		if err = a.updateOfflineToken(); err != nil {
			return
		}
	}
	if form == nil {
		form = core.NewForm()
	}
	form.WithInt("uid", a.ui.UserId).
		WithString("sign", a.ot.Sign).
		WithInt64("time", a.ot.Time)
	err = a.hc.JsonApi(url, nil, form, result)
	// TODO: handle token expired error.
	return
}

/*
Get some of offline tasks.

The upstream API returns at most 30 tasks for one request, caller need pass a cursor to
receive the cursor information, and use it to get remain tasks.

The cursor should be created by OfflineCursor(), DO NOT pass it as nil.
*/
func (a *Agent) OfflineList(cursor Cursor) (tasks []*OfflineTask, err error) {
	oc, ok := cursor.(*offlineCursor)
	if !ok {
		return nil, errOfflineCursorInvalid
	}
	form := core.NewForm().WithInt("page", oc.page)
	result := &internal.OfflineListResult{}
	err = a.callOfflineApi(apiOfflineList, form, result)
	if err == nil && result.IsFailed() {
		err = internal.MakeOfflineError(result.ErrorCode, result.ErrorMsg)
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
	// Update cursor
	oc.used, oc.total = true, result.Count
	oc.page, oc.pageCount = result.Page, result.PageCount
	return
}

// Add one or more offline tasks by URL.
func (a *Agent) OfflineAdd(url ...string) (err error) {
	form, isSingle := core.NewForm(), len(url) == 1
	if isSingle {
		form.WithString("url", url[0])
		result := &internal.OfflineAddUrlResult{}
		err = a.callOfflineApi(apiOfflineAddUrl, form, result)
	} else {
		form.WithStrings("url", url)
		result := &internal.OfflineAddUrlsResult{}
		err = a.callOfflineApi(apiOfflineAddUrls, form, result)
	}
	// TODO: return add result
	return
}

// Delete some offline tasks.
// if "deleteFile" is true, the downloaded files will be deleted.
func (a *Agent) OfflineDelete(deleteFile bool, hash ...string) (err error) {
	form := core.NewForm().WithStrings("hash", hash)
	if deleteFile {
		form.WithInt("flag", 1)
	}
	result := &internal.OfflineBasicResult{}
	err = a.callOfflineApi(apiOfflineDelete, form, result)
	if err == nil && !result.State {
		err = internal.MakeOfflineError(result.ErrorCode, result.ErrorMsg)
	}
	return
}

/*
Clear specific type of offline tasks.

The "flag" parameter indicates which type of tasks will be clear, you can pass nil to
clear the done tasks but keep the downloaded files.
*/
func (a *Agent) OfflineClear(flag *OfflineClearFlag) (err error) {
	if flag == nil {
		flag = (&OfflineClearFlag{}).Done(false)
	}
	form := core.NewForm().
		WithInt("flag", flag.flag)
	result := &internal.OfflineBasicResult{}
	err = a.callOfflineApi(apiOfflineClear, form, result)
	if err == nil && !result.State {
		err = internal.MakeOfflineError(result.ErrorCode, result.ErrorMsg)
	}
	return
}
