package api

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/api/base"
)

type _DownloadInfo struct {
	FileName string      `json:"file_name"`
	FileSize json.Number `json:"file_size"`
	PickCode string      `json:"pick_code"`
	Url      struct {
		Url    string `json:"url"`
		Client int    `json:"client"`
		OssId  string `json:"oss_id"`
	} `json:"url"`
}

type _DownloadData map[string]*_DownloadInfo

type DownloadSpec struct {
	base.M115ApiSpec[_DownloadData]
}

func (s *DownloadSpec) Init(pickcode string) *DownloadSpec {
	s.M115ApiSpec.Init("https://proapi.115.com/app/chrome/downurl")
	s.QuerySetNow("t")
	s.ParamSet("pickcode", pickcode)
	return s
}
