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

	offlineClearFlagMin = OfflineClearDone
	offlineClearFlagMax = OfflineClearAllAndDelete
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

func (t *OfflineTask) IsRunning() bool {
	return t.Status == 1
}

func (t *OfflineTask) IsDone() bool {
	return t.Status == 2
}

func (t *OfflineTask) IsFailed() bool {
	return t.Status == -1
}

func (a *Agent) offlineInitToken() (err error) {
	qs := web.Params{}.WithNow("_")
	resp := &webapi.OfflineSpaceResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiOfflineSpace, qs, nil, resp); err != nil {
		return
	}
	a.ot.Time = resp.Time
	a.ot.Sign = resp.Sign
	return nil
}

func (a *Agent) offlineCallApi(url string, form web.Params, resp web.ApiResp) (err error) {
	if a.ot.Time == 0 {
		if err = a.offlineInitToken(); err != nil {
			return
		}
	}
	if form == nil {
		form = web.Params{}
	}
	form.WithInt("uid", a.uid).
		WithInt64("time", a.ot.Time).
		With("sign", a.ot.Sign)
	return a.wc.CallJsonApi(url, nil, form, resp)
}

type OfflineIterator interface {
	Next() error
	Get(*OfflineTask) error
}

type implOfflineIterator struct {
	a *Agent
	// Page index
	pi int
	// Page count
	pc int
	// Cached tasks
	ts []*webapi.OfflineTask
	// Task index
	ti int
	// Task count
	tc int
}

func (i *implOfflineIterator) Next() (err error) {
	i.ti += 1
	// If we reach the last record?
	if i.ti == i.tc && i.pi == i.pc {
		return webapi.ErrReachEnd
	}
	// Is cache available?
	if i.ti < i.tc {
		return nil
	}
	// Fetch next page
	if err = i.a.offlineGetTasks(i.pi+1, i); err == nil {
		if i.tc == 0 {
			err = webapi.ErrReachEnd
		}
	}
	return
}

func (i *implOfflineIterator) Get(task *OfflineTask) (err error) {
	if i.ti >= i.tc {
		return webapi.ErrReachEnd
	}
	t := i.ts[i.ti]
	task.InfoHash = t.InfoHash
	task.Name = t.Name
	task.Size = t.Size
	task.Url = t.Url
	task.Status = t.Status
	task.Percent = t.Percent
	task.FileId = t.FileId
	return nil
}

// OfflineIterate returns an iterator for travelling offline tasks. it will
// return io.EOF error when there are no tasks.
func (a *Agent) OfflineIterate() (it OfflineIterator, err error) {
	impl := &implOfflineIterator{
		a: a,
	}
	if err = a.offlineGetTasks(1, impl); err == nil {
		if impl.tc == 0 {
			err = webapi.ErrReachEnd
		} else {
			it = impl
		}
	}
	return
}

func (a Agent) offlineGetTasks(page int, it *implOfflineIterator) (err error) {
	form := web.Params{}.
		WithInt("page", page)
	resp := &webapi.OfflineListResponse{}
	if err = a.offlineCallApi(webapi.ApiOfflineList, form, resp); err != nil {
		return
	}
	it.pc, it.pi = resp.PageCount, resp.PageIndex
	it.ts, it.tc, it.ti = resp.Tasks, len(resp.Tasks), 0
	return
}

// OfflineAdd adds an offline task with url, and saves the downloaded files at
// directory whose ID is dirId.
// You can pass empty string as dirId, to save the downloaded files at default
// directory.
func (a *Agent) OfflineAdd(url string, dirId string) (err error) {
	form := web.Params{}.
		With("url", url)
	if dirId != "" {
		form.With("wp_path_id", dirId)
	}
	resp := &webapi.OfflineAddUrlResponse{}
	return a.offlineCallApi(webapi.ApiOfflineAddUrl, form, resp)
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
	return a.offlineCallApi(
		webapi.ApiOfflineDelete, form, &webapi.OfflineBasicResponse{})
}

// OfflineClear clears tasks which is in specific status.
func (a *Agent) OfflineClear(flag OfflineClearFlag) (err error) {
	if flag < offlineClearFlagMin || flag > offlineClearFlagMax {
		flag = OfflineClearDone
	}
	form := web.Params{}.
		WithInt("flag", int(flag))
	return a.offlineCallApi(
		webapi.ApiOfflineClear, form, &webapi.OfflineBasicResponse{})
}
