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
	qs := protocol.Params{}.
		WithInt("page", oi.pi)
	resp := &webapi.OfflineListResponse{}
	if err = a.pc.CallSecretJsonApi(webapi.ApiTaskList, qs, nil, resp, 0); err != nil {
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

type OfflineAddResult map[string]*webapi.OfflineAddUrlResponse

func (r OfflineAddResult) Get(url string, task *OfflineTask) (err error) {
	if resp, ok := r[url]; !ok {
		return webapi.ErrNotExist
	} else if task != nil {
		task.InfoHash = resp.InfoHash
		task.Name = resp.Name
		task.Url = resp.Url
	}
	return
}

// OfflineAddUrl adds offline tasks from urls, this function calls 115 PC API
// which does not require CAPTCHA after you add a lot of tasks.
func (a *Agent) OfflineAddUrl(dirId string, urls []string, result OfflineAddResult) (err error) {
	// Prepare results buffer
	if urlCount := len(urls); urlCount == 0 {
		err = webapi.ErrEmptyList
		return
	}
	// Prepare request data
	params := protocol.Params{}.
		With("ac", "add_task_urls").
		With("app_ver", a.uh.AppVersion()).
		With("uid", a.uh.UserId()).
		WithArray("url", urls)
	if dirId != "" {
		params.With("wp_path_id", dirId)
	}
	data, err := json.Marshal(params)
	if err != nil {
		return
	}
	// M115 encoding
	key := m115.GenerateKey()
	form := protocol.Params{}.With("data", m115.Encode(data, key)).ToForm()
	mr := &webapi.M115Response{}
	if err = a.pc.CallJsonApi(webapi.ApiTaskAddUrls, nil, form, mr); err != nil {
		return
	}
	if data, err = m115.Decode(mr.Data, key); err != nil {
		return
	}
	resp := &webapi.OfflineAddUrlsResponse{}
	if err = json.Unmarshal(data, resp); err == nil && result != nil {
		for i, r := range resp.Result {
			result[urls[i]] = r
		}
	}
	return
}

// OfflineDelete deletes tasks.
func (a *Agent) OfflineDelete(deleteFiles bool, hashes []string) (err error) {
	if len(hashes) == 0 {
		return
	}
	form := protocol.Params{}.
		WithArray("hash", hashes)
	if deleteFiles {
		form.With("flag", "1")
	} else {
		form.With("flag", "0")
	}
	return a.pc.CallSecretJsonApi(
		webapi.ApiTaskDelete, nil, form.ToForm(),
		&webapi.OfflineBasicResponse{}, 0)
}

// OfflineClear clears tasks which is in specific status.
func (a *Agent) OfflineClear(flag OfflineClearFlag) (err error) {
	if flag < offlineClearFlagMin || flag > offlineClearFlagMax {
		flag = OfflineClearDone
	}
	form := protocol.Params{}.
		WithInt("flag", int(flag)).
		ToForm()
	return a.pc.CallSecretJsonApi(
		webapi.ApiTaskClear, nil, form,
		&webapi.OfflineBasicResponse{}, 0)
}
