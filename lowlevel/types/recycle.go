package types

import "github.com/deadblue/elevengo/internal/util"

type RecycleBinItem struct {
	FileId     string         `json:"id"`
	FileName   string         `json:"file_name"`
	FileSize   util.IntNumber `json:"file_size"`
	ParentId   string         `json:"cid"`
	ParentName string         `json:"parent_name"`
	DeleteTime util.IntNumber `json:"dtime"`
}

type RecycleBinListResult struct {
	Count int
	Item  []*RecycleBinItem
}
