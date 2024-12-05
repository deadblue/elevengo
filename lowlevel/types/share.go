package types

import "github.com/deadblue/elevengo/internal/util"

const (
	ShareStateAuditing   = 0
	ShareStateAccepted   = 1
	ShareStateRejected1  = 2
	ShareStateRejected2  = 3
	ShareStateCanceled   = 4
	ShareStateDeleted    = 5
	ShareStateRejected3  = 6
	ShareStateExpired    = 7
	ShareStateGenerating = 8
	ShareStateFailed     = 9
	ShareStateRejected4  = 11
)

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

type ShareFileInfo struct {
	FileId          string
	IsDir           bool
	Name            string
	Size            int64
	Sha1            string
	CreateTime      int64
	IsVideo         bool
	VideoDefinition int
	MediaDuration   int
}

type ShareSnapResult struct {
	SnapId       string
	UserId       int
	ShareTitle   string
	ShareState   int
	ReceiveCount int
	CreateTime   int64
	ExpireTime   int64

	TotalSize int64
	FileCount int
	Files     []*ShareFileInfo
}
