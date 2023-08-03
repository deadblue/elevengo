package api

import "github.com/deadblue/elevengo/internal/api/base"

type ImageGetResult struct {
	FileName string `json:"file_name"`
	FileSha1 string `json:"file_sha1"`
	Pickcode string `json:"pick_code"`

	SourceUrl string   `json:"source_url"`
	OriginUrl string   `json:"origin_url"`
	ViewUrl   string   `json:"url"`
	ThumbUrls []string `json:"all_url"`
}

type ImageGetSpec struct {
	base.JsonApiSpec[ImageGetResult, base.StandardResp]
}

func (s *ImageGetSpec) Init(pickcode string) *ImageGetSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/image")
	s.QuerySet("pickcode", pickcode)
	s.QuerySetNow("_")
	return s
}
