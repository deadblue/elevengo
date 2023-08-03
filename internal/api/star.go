package api

import "github.com/deadblue/elevengo/internal/api/base"

type FileStarSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *FileStarSpec) Init(fileId string, star bool) *FileStarSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/files/star")
	s.FormSet("file_id", fileId)
	if star {
		s.FormSet("star", "1")
	} else {
		s.FormSet("star", "0")
	}
	return s
}
