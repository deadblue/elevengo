package types

import "github.com/deadblue/elevengo/internal/util"

type ShareInfo struct {
	ShareCode string `json:"share_code"`

	ShareState    util.IntNumber `json:"share_state"`
	ShareTitle    string         `json:"share_title"`
	ShareUrl      string         `json:"share_url"`
	ShareDuration util.IntNumber `json:"share_ex_time"`
	ReceiveCode   string         `json:"receive_code"`

	ReceiveCount util.IntNumber `json:"receive_count"`

	FileCount   int            `json:"file_count"`
	FolderCount int            `json:"folder_count"`
	TotalSize   util.IntNumber `json:"total_size"`
}

type ShareListResult struct {
	Count int
	Items []*ShareInfo
}
