package elevengo

import (
	"github.com/deadblue/elevengo/internal/web"
	"github.com/deadblue/elevengo/internal/webapi"
)

type OfflineClearFlag int

const (
	OfflineClearDone OfflineClearFlag = iota
	OfflineClearAll
	OfflineClearFailed
	OfflineClearRunning
	OfflineClearDoneAndDelete
	OfflineClearAllAndDelete
)

// OfflineTask describe an offline downloading task.
type OfflineTask struct {
	InfoHash string
	Name     string
	Size     int64
	Status   int
	Percent  float64
	Url      string
	FileId   string
}

func (t *OfflineTask) from(task *webapi.OfflineTask) {
	t.InfoHash = task.InfoHash
	t.Name = task.Name
	t.Size = task.Size
	t.Status = task.Status
	t.Percent = task.Percent
	t.Url = task.Url
	t.FileId = task.FileId
}

func (t *OfflineTask) IsRunning() bool {
	return t.Status == 1
}

func (t *OfflineTask) IsDone() bool {
	return t.Status == 2
}

func (t *OfflineTask) IsFailed() bool {
	return t.Status == -1
}

type OfflineCursor struct {
	page int
}

func (a *Agent) offlineUpdateToken() (err error) {
	qs := web.Params{}.WithNow("_")
	resp := &webapi.OfflineSpaceResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiOfflineSpace, qs, nil, resp); err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	a.ot.Time = resp.Time
	a.ot.Sign = resp.Sign
	return nil
}

func (a *Agent) offlineCallApi(url string, form web.Params, resp interface{}) (err error) {
	if a.ot.Time == 0 {
		if err = a.offlineUpdateToken(); err != nil {
			return
		}
	}
	if form == nil {
		form = web.Params{}
	}
	form.WithInt("uid", a.user.Id).
		WithInt64("time", a.ot.Time).
		With("sign", a.ot.Sign)
	return a.wc.CallJsonApi(url, nil, form, resp)
}

// OfflineList lists offline tasks
func (a *Agent) OfflineList() (err error) {
	form := web.Params{}.
		WithInt("page", 1)
	resp := &webapi.OfflineListResponse{}
	if err = a.offlineCallApi(webapi.ApiOfflineList, form, resp); err != nil {
		err = resp.Err()
	}
	if err != nil {
		return
	}
	// TODO: How we return
	return
}

// OfflineAdd adds an offline task with url.
func (a *Agent) OfflineAdd(url string, dirId string) (err error) {
	form := web.Params{}.
		With("url", url)
	if dirId != "" {
		form.With("wp_path_id", dirId)
	}
	resp := &webapi.OfflineAddUrlResponse{}
	if err = a.offlineCallApi(webapi.ApiOfflineAddUrl, form, resp); err != nil {
		err = resp.Err()
	}
	return
}

// OfflineDelete deletes tasks.
func (a *Agent) OfflineDelete(deleteFiles bool, hashes ...string) (err error) {
	if len(hashes) == 0 {
		return
	}
	form := web.Params{}.WithArray("hash", hashes)
	if deleteFiles {
		form.With("flag", "1")
	}
	resp := &webapi.OfflineBasicResponse{}
	if err = a.offlineCallApi(webapi.ApiOfflineDelete, form, resp); err != nil {
		err = resp.Err()
	}
	return
}

// OfflineClear clears tasks which is in specific status.
func (a *Agent) OfflineClear(flag OfflineClearFlag) (err error) {
	form := web.Params{}.
		WithInt("flag", int(flag))
	resp := &webapi.OfflineBasicResponse{}
	if err = a.offlineCallApi(webapi.ApiOfflineClear, form, resp); err != nil {
		err = resp.Err()
	}
	return
}
