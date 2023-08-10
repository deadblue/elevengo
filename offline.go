package elevengo

import (
	"github.com/deadblue/elevengo/internal/api"
	"github.com/deadblue/elevengo/internal/api/errors"
	"github.com/deadblue/elevengo/option"
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
	tasks []*api.OfflineTask
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
		return errors.ErrReachEnd
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
		return errors.ErrReachEnd
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
	spec := (&api.OfflineListSpec{}).Init(oi.pi)
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	result := spec.Result
	oi.pi = result.PageIndex
	oi.pc = result.PageCount
	oi.ps = result.PageSize
	oi.index, oi.size = 0, len(result.Tasks)
	if oi.size == 0 {
		err = errors.ErrReachEnd
	} else {
		oi.tasks = make([]*api.OfflineTask, 0, oi.size)
		oi.tasks = append(oi.tasks, result.Tasks...)
	}
	oi.count = result.TaskCount
	return
}

// OfflineDelete deletes tasks.
func (a *Agent) OfflineDelete(hashes []string, opts ...option.OfflineDeleteOption) (err error) {
	if len(hashes) == 0 {
		return
	}
	// Apply options
	deleteFiles := false
	for _, opt := range opts {
		switch opt := opt.(type) {
		case option.OfflineDeleteFilesOfTasks:
			deleteFiles = bool(opt)
		}
	}
	// Call API
	spec := (&api.OfflineDeleteSpec{}).Init(hashes, deleteFiles)
	return a.pc.ExecuteApi(spec)
}

// OfflineClear clears tasks which is in specific status.
func (a *Agent) OfflineClear(flag OfflineClearFlag) (err error) {
	if flag < offlineClearFlagMin || flag > offlineClearFlagMax {
		flag = OfflineClearDone
	}
	spec := (&api.OfflineClearSpec{}).Init(int(flag))
	return a.pc.ExecuteApi(spec)
}

// OfflineAddUrl adds offline tasks by download URLs.
// It returns an info hash list related to the given urls, the info hash will
// be empty if the related URL is invalid.
//
// You can use options to change the download directory:
//
//	agent := Default()
//	agent.CredentialImport(&Credential{UID: "", CID: "", SEID: ""})
//	hashes, err := agent.OfflineAddUrl([]string{
//		"https://foo.bar/file.zip",
//		"magent:?xt=urn:btih:111222",
//		"ed2k://|file|name|size|md4|",
//	}, option.OfflineSaveDownloadedFileTo("dirId"))
func (a *Agent) OfflineAddUrl(urls []string, opts ...option.OfflineAddOption) (hashes []string, err error) {
	// Prepare results buffer
	if urlCount := len(urls); urlCount == 0 {
		return
	} else {
		hashes = make([]string, urlCount)
	}
	// Apply options
	saveDirId := ""
	for _, opt := range opts {
		switch opt := opt.(type) {
		case option.OfflineSaveDownloadedFileTo:
			saveDirId = string(opt)
		}
	}
	// Call API
	spec := (&api.OfflineAddUrlsSpec{}).Init(
		a.uh.UserId, a.uh.AppVer, urls, saveDirId,
	)
	if err = a.pc.ExecuteApi(spec); err == nil {
		for i, task := range spec.Result {
			if task != nil {
				hashes[i] = task.InfoHash
			}
		}
	}
	return
}
