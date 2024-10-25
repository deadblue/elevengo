package elevengo

import (
	"context"
	"iter"

	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/types"
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

func (t *OfflineTask) from(ti *types.TaskInfo) *OfflineTask {
	t.InfoHash = ti.InfoHash
	t.Name = ti.Name
	t.Size = ti.Size
	t.Status = ti.Status
	t.Percent = ti.Percent
	t.Url = ti.Url
	t.FileId = ti.FileId
	return t
}

type offlineIterator struct {
	llc    client.Client
	page   int
	result *types.OfflineListResult
}

func (i *offlineIterator) update() (err error) {
	if i.result != nil && i.page > i.result.PageCount {
		return errNoMoreItems
	}
	spec := (&api.OfflineListSpec{}).Init(i.page)
	if err = i.llc.CallApi(spec, context.Background()); err == nil {
		i.result = &spec.Result
		i.page += 1
	}
	return
}

func (i *offlineIterator) Count() int {
	if i.result == nil {
		return 0
	}
	return i.result.TaskCount
}

func (i *offlineIterator) Items() iter.Seq2[int, *OfflineTask] {
	return func(yield func(int, *OfflineTask) bool) {
		for index := 0; ; {
			for _, ti := range i.result.Tasks {
				if stop := !yield(index, (&OfflineTask{}).from(ti)); stop {
					return
				}
				index += 1
			}
			if err := i.update(); err != nil {
				break
			}
		}
	}
}

// OfflineIterate returns an iterator to access all offline tasks.
func (a *Agent) OfflineIterate() (it Iterator[OfflineTask], err error) {
	oi := &offlineIterator{
		llc:  a.llc,
		page: 1,
	}
	if err = oi.update(); err == nil {
		it = oi
	}
	return
}

// OfflineDelete deletes tasks.
func (a *Agent) OfflineDelete(hashes []string, options ...*option.OfflineDeleteOptions) (err error) {
	if len(hashes) == 0 {
		return
	}
	// Apply options
	deleteFiles := false
	if opts := util.NotNull(options...); opts != nil {
		deleteFiles = opts.DeleteFiles
	}
	// Call API
	spec := (&api.OfflineDeleteSpec{}).Init(hashes, deleteFiles)
	return a.llc.CallApi(spec, context.Background())
}

// OfflineClear clears tasks which is in specific status.
func (a *Agent) OfflineClear(flag OfflineClearFlag) (err error) {
	if flag < offlineClearFlagMin || flag > offlineClearFlagMax {
		flag = OfflineClearDone
	}
	spec := (&api.OfflineClearSpec{}).Init(int(flag))
	return a.llc.CallApi(spec, context.Background())
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
func (a *Agent) OfflineAddUrl(urls []string, options ...*option.OfflineAddOptions) (hashes []string, err error) {
	// Prepare results buffer
	if urlCount := len(urls); urlCount == 0 {
		return
	} else {
		hashes = make([]string, urlCount)
	}
	// Apply options
	saveDirId := ""
	if opts := util.NotNull(options...); opts != nil {
		saveDirId = opts.SaveDirId
	}
	// Call API
	spec := (&api.OfflineAddUrlsSpec{}).Init(urls, saveDirId, &a.common)
	if err = a.llc.CallApi(spec, context.Background()); err == nil {
		for i, task := range spec.Result {
			if task != nil {
				hashes[i] = task.InfoHash
			}
		}
	}
	return
}
