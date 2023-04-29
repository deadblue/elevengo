package elevengo

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/crypto/m115"
	"github.com/deadblue/elevengo/internal/protocol"
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
	qs := protocol.Params{}.WithNow("_")
	resp := &webapi.OfflineSpaceResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiOfflineSpace, qs, nil, resp); err != nil {
		return
	}
	a.ot.Time = resp.Time
	a.ot.Sign = resp.Sign
	return nil
}

func (a *Agent) offlineCallApi(url string, params protocol.Params, resp protocol.ApiResp) (err error) {
	if a.ot.Time == 0 {
		if err = a.offlineInitToken(); err != nil {
			return
		}
	}
	if params == nil {
		params = protocol.Params{}
	}
	params.WithInt("uid", a.uid).
		WithInt64("time", a.ot.Time).
		With("sign", a.ot.Sign)
	return a.pc.CallJsonApi(url, nil, params.ToForm(), resp)
}

type offlineIterator struct {
	// Total task count
	count int
	// Page index
	pi int
	// Page count
	pc int
	// Page size
	ps int

	// Cached tasks
	tasks []*webapi.OfflineTask
	// Task index
	index int
	// Task size
	size int

	// Update function
	uf func(*offlineIterator) error
}

func (i *offlineIterator) Next() (err error) {
	if i.index += 1; i.index < i.size {
		return nil
	}
	if i.pi >= i.pc {
		return webapi.ErrReachEnd
	}
	// Fetch next page
	i.pi += 1
	return i.uf(i)
}

func (i *offlineIterator) Index() int {
	return (i.pi-1)*i.ps + i.index
}

func (i *offlineIterator) Get(task *OfflineTask) (err error) {
	if i.index >= i.size {
		return webapi.ErrReachEnd
	}
	t := i.tasks[i.index]
	task.InfoHash = t.InfoHash
	task.Name = t.Name
	task.Size = t.Size
	task.Url = t.Url
	task.Status = t.Status
	task.Percent = t.Percent
	task.FileId = t.FileId
	return nil
}

func (i *offlineIterator) Count() int {
	return i.count
}

// OfflineIterate returns an iterator for travelling offline tasks, it will
// return an error if there are no tasks.
func (a *Agent) OfflineIterate() (it Iterator[OfflineTask], err error) {
	oi := &offlineIterator{
		pi: 1,
		uf: a.offlineIterateInternal,
	}
	if err = a.offlineIterateInternal(oi); err == nil {
		it = oi
	}
	return
}

func (a *Agent) offlineIterateInternal(oi *offlineIterator) (err error) {
	form := protocol.Params{}.
		WithInt("page", oi.pi)
	resp := &webapi.OfflineListResponse{}
	if err = a.offlineCallApi(webapi.ApiOfflineList, form, resp); err != nil {
		return
	}
	oi.pi = resp.PageIndex
	oi.pc = resp.PageCount
	oi.ps = resp.PageSize
	oi.index, oi.size = 0, len(resp.Tasks)
	if oi.size == 0 {
		err = webapi.ErrReachEnd
	} else {
		oi.tasks = make([]*webapi.OfflineTask, 0, oi.size)
		oi.tasks = append(oi.tasks, resp.Tasks...)
	}
	oi.count = resp.TaskCount
	return
}

type OfflineAddResult struct {
	InfoHash string
	Name     string
	Error    error
}

func (r *OfflineAddResult) IsExist() bool {
	return r.Error == webapi.ErrOfflineTaskExisted
}

// Deprecated: Please use `OfflineAddUrl` instead.
// 
// OfflineAdd adds an offline task with url, and saves the downloaded files at
// directory whose ID is dirId.
// You can pass empty string as dirId, to save the downloaded files at default
// directory.
func (a *Agent) OfflineAdd(url string, dirId string) (result OfflineAddResult) {
	form := protocol.Params{}.
		With("url", url)
	if dirId != "" {
		form.With("wp_path_id", dirId)
	}
	resp := &webapi.OfflineAddUrlResponse{}
	result.Error = a.offlineCallApi(webapi.ApiOfflineAddUrl, form, resp)
	result.InfoHash, result.Name = resp.InfoHash, resp.Name
	return
}

// Deprecated: Please use `OfflineAddUrl` instead.
//
// OfflineBatchAdd adds many offline tasks in one request.
func (a *Agent) OfflineBatchAdd(urls []string, dirId string) (results []OfflineAddResult, err error) {
	if urlCount := len(urls); urlCount == 0 {
		err = webapi.ErrEmptyList
		return
	} else {
		results = make([]OfflineAddResult, urlCount)
	}

	form := protocol.Params{}.
		WithArray("url", urls)
	if dirId != "" {
		form.With("wp_path_id", dirId)
	}
	resp := &webapi.OfflineAddUrlsResponse{}
	if err = a.offlineCallApi(webapi.ApiOfflineAddUrls, form, resp); err != nil {
		return
	}
	for i, result := range resp.Result {
		results[i].InfoHash = result.InfoHash
		results[i].Name = result.Name
		results[i].Error = result.Err()
	}
	return
}

// OfflineAddUrl adds offline tasks from urls, this API calls 115 PC API 
// which (may) not require captcha after you add a lot of tasks.
func (a *Agent) OfflineAddUrl(urls ...string) (results []OfflineAddResult, err error) {
	// Prepare results buffer
	if urlCount := len(urls); urlCount == 0 {
		err = webapi.ErrEmptyList
		return
	} else {
		results = make([]OfflineAddResult, urlCount)
	}
	// Prepare request data
	params := protocol.Params{}.
		With("ac", "add_task_urls").
		With("app_ver", a.uh.AppVersion()).
		WithInt("uid", a.uid).
		WithArray("url", urls)
	// if dirId != "" {
	// 	params.With("savepath", dirId)
	// }
	data ,err := json.Marshal(params)
	if err != nil {
		return 
	}
	key := m115.GenerateKey()
	form := protocol.Params{}.With("data", m115.Encode(data, key)).ToForm()
	mr := &webapi.M115Response{}
	if err = a.pc.CallJsonApi(webapi.ApiOfflineAddUrlsNew, nil, form, mr); err != nil {
		return
	}
	if data, err = m115.Decode(mr.Data, key); err != nil {
		return
	}
	resp := &webapi.OfflineAddUrlsResponse{}
	if err = json.Unmarshal(data, resp); err != nil {
		return
	}
	for i, result := range resp.Result {
		results[i].InfoHash = result.InfoHash
		results[i].Name = result.Name
		results[i].Error = result.Err()
	}
	return
}

// OfflineDelete deletes tasks.
func (a *Agent) OfflineDelete(deleteFiles bool, hashes ...string) (err error) {
	if len(hashes) == 0 {
		return
	}
	form := protocol.Params{}.WithArray("hash", hashes)
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
	form := protocol.Params{}.
		WithInt("flag", int(flag))
	return a.offlineCallApi(
		webapi.ApiOfflineClear, form, &webapi.OfflineBasicResponse{})
}
