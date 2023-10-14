package types

type ShortcutInfo struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	Sort     string `json:"sort"`
}

type ShortcutListResult struct {
	List []*ShortcutInfo `json:"list"`
}
