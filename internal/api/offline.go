package api

import (
	"fmt"

	"github.com/deadblue/elevengo/internal/api/base"
	"github.com/deadblue/elevengo/internal/api/errors"
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

type _OfflineListData struct {
	PageIndex int
	PageCount int
	PageSize  int

	TaskCount int
	Tasks     []*OfflineTask
}

type _OfflineListResp struct {
	base.BasicResp

	PageIndex int `json:"page"`
	PageCount int `json:"page_count"`
	PageSize  int `json:"page_row"`

	QuotaTotal  int `json:"total"`
	QuotaRemain int `json:"quota"`

	TaskCount int            `json:"count"`
	Tasks     []*OfflineTask `json:"tasks"`
}

func (r *_OfflineListResp) Extract(v any) (err error) {
	if ptr, ok := v.(*_OfflineListData); !ok {
		return errors.ErrUnsupportedData
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
	base.JsonApiSpec[_OfflineListResp, _OfflineListData]
}

func (s *OfflineListSpec) Init(page int) *OfflineListSpec {
	s.JsonApiSpec.Init("https://lixian.115.com/lixian/?ct=lixian&ac=task_lists")
	s.QuerySetInt("page", page)
	return s
}

func (s *OfflineListSpec) SetPage(page int) {
	s.QuerySetInt("page", page)
}

type _OfflineAddResult struct {
	State   bool `json:"state"`
	ErrNum  int  `json:"errno"`
	ErrCode int  `json:"errcode"`

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

type OfflineAddUrlsSpec struct {
	base.M115ApiSpec[_OfflineAddUrlsData]
}

func (s *OfflineAddUrlsSpec) Init(userId, appVer string, urls []string) *OfflineAddUrlsSpec {
	s.M115ApiSpec.Init("https://lixian.115.com/lixianssp/?ac=add_task_urls")
	s.M115ApiSpec.EnableCrypto()
	s.ParamSetAll(map[string]string{
		"ac":      "add_task_urls",
		"app_ver": appVer,
		"uid":     userId,
	})
	for i, url := range urls {
		key := fmt.Sprintf("url[%d]", i)
		s.ParamSet(key, url)
	}
	return s
}

type OfflineDeleteSpec struct {
	base.JsonApiSpec[base.StandardResp, base.VoidData]
}

func (s *OfflineDeleteSpec) Init(hashes []string, deleteFiles bool) *OfflineDeleteSpec {
	s.JsonApiSpec.Init("https://lixian.115.com/lixian/?ct=lixian&ac=task_del")
	if deleteFiles {
		s.FormSet("flag", "1")
	} else {
		s.FormSet("flag", "0")
	}
	return s
}

type OfflineClearSpec struct {
	base.JsonApiSpec[base.BasicResp, base.VoidData]
}

func (s *OfflineClearSpec) Init(flag int) *OfflineClearSpec {
	s.JsonApiSpec.Init("https://lixian.115.com/lixian/?ct=lixian&ac=task_clear")
	s.FormSetInt("flag", flag)
	return s
}
