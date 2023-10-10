package api

import (
	"encoding/json"
	"fmt"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/errors"
)

type OfflineTask struct {
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Url      string `json:"url"`
	AddTime  int64  `json:"add_time"`

	Status     int     `json:"status"`
	Percent    float64 `json:"percentDone"`
	UpdateTime int64   `json:"last_update"`

	FileId string `json:"file_id"`
	DirId  string `json:"wp_path_id"`
}

type OfflineListResult struct {
	PageIndex int
	PageCount int
	PageSize  int

	TaskCount int
	Tasks     []*OfflineTask
}

//lint:ignore U1000 This type is used in generic.
type _OfflineListResp struct {
	protocol.BasicResp

	PageIndex int `json:"page"`
	PageCount int `json:"page_count"`
	PageSize  int `json:"page_row"`

	QuotaTotal  int `json:"total"`
	QuotaRemain int `json:"quota"`

	TaskCount int            `json:"count"`
	Tasks     []*OfflineTask `json:"tasks"`
}

func (r *_OfflineListResp) Extract(v any) (err error) {
	if ptr, ok := v.(*OfflineListResult); !ok {
		return errors.ErrUnsupportedResult
	} else {
		ptr.PageIndex = r.PageIndex
		ptr.PageCount = r.PageCount
		ptr.PageSize = r.PageSize
		ptr.TaskCount = r.TaskCount
		ptr.Tasks = append(ptr.Tasks, r.Tasks...)
	}
	return nil
}

type OfflineListSpec struct {
	_JsonApiSpec[OfflineListResult, _OfflineListResp]
}

func (s *OfflineListSpec) Init(page int) *OfflineListSpec {
	s._JsonApiSpec.Init("https://lixian.115.com/lixian/?ct=lixian&ac=task_lists")
	s.query.SetInt("page", page)
	return s
}

type OfflineDeleteSpec struct {
	_VoidApiSpec
}

func (s *OfflineDeleteSpec) Init(hashes []string, deleteFiles bool) *OfflineDeleteSpec {
	s._VoidApiSpec.Init("https://lixian.115.com/lixian/?ct=lixian&ac=task_del")
	for index, hash := range hashes {
		key := fmt.Sprintf("hash[%d]", index)
		s.form.Set(key, hash)
	}
	if deleteFiles {
		s.form.Set("flag", "1")
	} else {
		s.form.Set("flag", "0")
	}
	return s
}

type OfflineClearSpec struct {
	_VoidApiSpec
}

func (s *OfflineClearSpec) Init(flag int) *OfflineClearSpec {
	s._VoidApiSpec.Init("https://lixian.115.com/lixian/?ct=lixian&ac=task_clear")
	s.form.SetInt("flag", flag)
	return s
}

type _OfflineAddResult struct {
	State   bool   `json:"state"`
	ErrNum  int    `json:"errno"`
	ErrCode int    `json:"errcode"`
	ErrType string `json:"errtype"`

	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

type _OfflineAddUrlsData struct {
	State   bool `json:"state"`
	ErrNum  int  `json:"errno"`
	ErrCode int  `json:"errcode"`

	Result []*_OfflineAddResult `json:"result"`
}

type OfflineAddUrlsResult []*OfflineTask

type OfflineAddUrlsSpec struct {
	_M115ApiSpec[OfflineAddUrlsResult]
}

func offlineAddUrlsResultExtractor(data []byte, result *OfflineAddUrlsResult) (err error) {
	obj := &_OfflineAddUrlsData{}
	if err = json.Unmarshal(data, obj); err != nil {
		return
	}
	tasks := make([]*OfflineTask, len(obj.Result))
	for i, r := range obj.Result {
		if r.State || r.ErrCode == errors.CodeOfflineTaskExists {
			tasks[i] = &OfflineTask{}
			tasks[i].InfoHash = r.InfoHash
			tasks[i].Name = r.Name
			tasks[i].Url = r.Url
		} else {
			tasks[i] = nil
		}
	}
	*result = tasks
	return
}

func (s *OfflineAddUrlsSpec) Init(userId, appVer string, urls []string, saveDirId string) *OfflineAddUrlsSpec {
	s._M115ApiSpec.Init(
		"https://lixian.115.com/lixianssp/?ac=add_task_urls",
		offlineAddUrlsResultExtractor,
	)
	s.crypto = true
	s.params.Set("uid", userId).
		Set("app_ver", appVer).
		Set("ac", "add_task_urls")
	for i, url := range urls {
		key := fmt.Sprintf("url[%d]", i)
		s.params.Set(key, url)
	}
	if saveDirId != "" {
		s.params.Set("wp_path_id", saveDirId)
	}
	return s
}
