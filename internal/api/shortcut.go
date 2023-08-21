package api

import "github.com/deadblue/elevengo/internal/api/base"

type ShortcutInfo struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	Sort     string `json:"sort"`
}

type ShortcutListResult struct {
	List []ShortcutInfo `json:"list"`
}

type ShortcutListSpec struct {
	base.JsonApiSpec[ShortcutListResult, base.StandardResp]
}

func (s *ShortcutListSpec) Init(fileId string) *ShortcutListSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/category/shortcut")
	return s
}

type ShortcutAddSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *ShortcutAddSpec) Init(fileId string) *ShortcutAddSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/category/shortcut")
	s.FormSet("file_id", fileId)
	s.FormSet("op", "add")
	return s
}

type ShortcutDeleteSpec struct {
	base.JsonApiSpec[base.VoidResult, base.BasicResp]
}

func (s *ShortcutDeleteSpec) Init(fileId string) *ShortcutDeleteSpec {
	s.JsonApiSpec.Init("https://webapi.115.com/category/shortcut")
	s.FormSet("file_id", fileId)
	s.FormSet("op", "delete")
	return s
}
