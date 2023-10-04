package api

import (
	"github.com/deadblue/elevengo/internal/apibase"
)

type ShortcutInfo struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	Sort     string `json:"sort"`
}

type ShortcutListResult struct {
	List []ShortcutInfo `json:"list"`
}

type ShortcutListSpec struct {
	apibase.StandardApiSpec[ShortcutListResult]
}

func (s *ShortcutListSpec) Init(fileId string) *ShortcutListSpec {
	s.StandardApiSpec.Init("https://webapi.115.com/category/shortcut")
	return s
}

type ShortcutAddSpec struct {
	apibase.VoidApiSpec
}

func (s *ShortcutAddSpec) Init(fileId string) *ShortcutAddSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/category/shortcut")
	s.FormSet("file_id", fileId)
	s.FormSet("op", "add")
	return s
}

type ShortcutDeleteSpec struct {
	apibase.VoidApiSpec
}

func (s *ShortcutDeleteSpec) Init(fileId string) *ShortcutDeleteSpec {
	s.VoidApiSpec.Init("https://webapi.115.com/category/shortcut")
	s.FormSet("file_id", fileId)
	s.FormSet("op", "delete")
	return s
}
