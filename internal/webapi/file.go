package webapi

import "encoding/json"

const (
	FileOrderByTime = "user_ptime"
	FileOrderByType = "file_type"
	FileOrderBySize = "file_size"
	FileOrderByName = "file_name"

	FileListLimit = 32
)

type FileInfo struct {
	AreaId     IntString `json:"aid"`
	CategoryId string    `json:"cid"`
	FileId     string    `json:"fid"`
	ParentId   string    `json:"pid"`

	Name     string      `json:"n"`
	Type     string      `json:"ico"`
	Size     StringInt64 `json:"s"`
	Sha1     string      `json:"sha"`
	PickCode string      `json:"pc"`

	IsStar StringInt    `json:"m"`
	Labels []*LabelInfo `json:"fl"`

	CreatedTime  string `json:"tp"`
	UpdatedTime  string `json:"te"`
	ModifiedTime string `json:"t"`

	// MediaDuration describes duration in seconds for audio / video.
	MediaDuration float64 `json:"play_long"`

	// Special fields for video
	IsVideo         int `json:"iv"`
	VideoDefinition int `json:"vdi"`
}

type FileParentInfo struct {
	FileId   IntString `json:"file_id"`
	FileName string    `json:"file_name"`
}

type FileDuplication struct {
	FileId   string      `json:"fid"`
	Name     string      `json:"fn"`
	Type     string      `json:"ico"`
	Size     StringInt64 `json:"fs"`
	Sha1     string      `json:"sha1"`
	PickCode string      `json:"pc"`

	UpdateTime json.Number `json:"te"`

	Parent FileParentInfo `json:"paths"`
}

type FileListResponse struct {
	BasicResponse

	AreaId     string    `json:"aid"`
	CategoryId IntString `json:"cid"`

	Count int    `json:"count"`
	Order string `json:"order"`
	IsAsc int    `json:"is_asc"`

	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
	PageSize int `json:"page_size"`
}

type FileStatResponse struct {
	BasicResponse

	FileName string    `json:"file_name"`
	PickCode string    `json:"pick_code"`
	Sha1     string    `json:"sha1"`
	IsFile   StringInt `json:"file_category"`

	CreateTime json.Number `json:"ptime"`
	UpdateTime json.Number `json:"utime"`
	AccessTime int64       `json:"open_time"`

	Paths []*FileParentInfo `json:"paths"`

	FileCount StringInt `json:"count"`
	DirCount  StringInt `json:"folder_count"`
}

func (r *FileStatResponse) Err() error {
	if r.PickCode != "" {
		return nil
	}
	return r.BasicResponse.Err()
}
