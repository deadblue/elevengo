package api

import (
	"fmt"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/lowlevel/types"
)

type OfflineListSpec struct {
	_JsonApiSpec[types.OfflineListResult, protocol.OfflineListResp]
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

type OfflineAddUrlsSpec struct {
	_M115ApiSpec[types.OfflineAddUrlsResult]
}

func (s *OfflineAddUrlsSpec) Init(userId, appVer string, urls []string, saveDirId string) *OfflineAddUrlsSpec {
	s._M115ApiSpec.Init("https://lixian.115.com/lixianssp/?ac=add_task_urls")
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
