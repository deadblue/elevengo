package api

import "github.com/deadblue/elevengo/lowlevel/types"

type ShortcutListSpec struct {
	_StandardApiSpec[types.ShortcutListResult]
}

func (s *ShortcutListSpec) Init(fileId string) *ShortcutListSpec {
	s._StandardApiSpec.Init("https://webapi.115.com/category/shortcut")
	return s
}

type ShortcutAddSpec struct {
	_VoidApiSpec
}

func (s *ShortcutAddSpec) Init(fileId string) *ShortcutAddSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/category/shortcut")
	s.form.Set("file_id", fileId).Set("op", "add")
	return s
}

type ShortcutDeleteSpec struct {
	_VoidApiSpec
}

func (s *ShortcutDeleteSpec) Init(fileId string) *ShortcutDeleteSpec {
	s._VoidApiSpec.Init("https://webapi.115.com/category/shortcut")
	s.form.Set("file_id", fileId).Set("op", "delete")
	return s
}
