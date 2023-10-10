package api

type ShortcutInfo struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	Sort     string `json:"sort"`
}

type ShortcutListResult struct {
	List []*ShortcutInfo `json:"list"`
}

type ShortcutListSpec struct {
	_StandardApiSpec[ShortcutListResult]
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
