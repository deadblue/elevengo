package types

import "github.com/deadblue/elevengo/internal/util"

type FileInfo struct {
	AreaId     util.IntNumber `json:"aid"`
	CategoryId string         `json:"cid"`
	FileId     string         `json:"fid"`
	ParentId   string         `json:"pid"`

	Name     string         `json:"n"`
	Type     string         `json:"ico"`
	Size     util.IntNumber `json:"s"`
	Sha1     string         `json:"sha"`
	PickCode string         `json:"pc"`

	IsStar util.Boolean `json:"m"`
	Labels []LabelInfo  `json:"fl"`

	CreatedTime  string `json:"tp"`
	UpdatedTime  string `json:"te"`
	ModifiedTime string `json:"t"`

	// MediaDuration describes duration in seconds for audio/video.
	MediaDuration float64 `json:"play_long"`

	// Special fields for video
	VideoFlag       int `json:"iv"`
	VideoDefinition int `json:"vdi"`
}

type FileListResult struct {
	DirId  string
	Offset int

	Count int
	Files []*FileInfo

	// Order settings
	Order string
	Asc   int
}

type FileGetResult []*FileInfo
