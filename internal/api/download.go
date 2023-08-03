package api

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/api/base"
)

type _DownloadUrlObject struct {
	Url    string `json:"url"`
	Client int    `json:"client"`
	OssId  string `json:"oss_id"`
}

type _DownloadUrl struct {
	Url string
}

func (u *_DownloadUrl) UnmarshalJSON(b []byte) (err error) {
	if len(b) > 0 && b[0] == '{' {
		obj := &_DownloadUrlObject{}
		if err = json.Unmarshal(b, obj); err == nil {
			u.Url = obj.Url
		}
	}
	return
}

type _DownloadInfo struct {
	FileName string       `json:"file_name"`
	FileSize json.Number  `json:"file_size"`
	PickCode string       `json:"pick_code"`
	Url      _DownloadUrl `json:"url"`
}

type DownloadResult map[string]*_DownloadInfo

type DownloadSpec struct {
	base.M115ApiSpec[DownloadResult]
}

func (s *DownloadSpec) Init(pickcode string) *DownloadSpec {
	s.M115ApiSpec.Init("https://proapi.115.com/app/chrome/downurl")
	s.QuerySetNow("t")
	s.ParamSet("pickcode", pickcode)
	return s
}
