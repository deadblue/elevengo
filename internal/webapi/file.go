package webapi

type FileInfo struct {
	AreaId     IntString `json:"aid"`
	CategoryId string    `json:"cid"`
	FileId     string    `json:"fid"`
	ParentId   string    `json:"pid"`

	Name     string      `json:"n"`
	Type     string      `json:"ico"`
	Size     StringInt64 `json:"s"`
	PickCode string      `json:"pc"`
	Sha1     string      `json:"sha"`

	CreateTime StringInt64 `json:"tp"`
	UpdateTime StringInt64 `json:"te"`

	// MediaDuration describes duration in seconds for audio / video.
	MediaDuration float64 `json:"play_long"`

	// Special fields for video
	IsVideo         int `json:"iv"`
	VideoDefinition int `json:"vdi"`
}

type DirInfo struct {
	FileId   IntString `json:"file_id"`
	FileName string    `json:"file_name"`
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

type FileSearchResponse struct {
	BasicResponse

	Count int    `json:"count"`
	Order string `json:"order"`
	IsAsc int    `json:"is_asc"`

	Offset   int `json:"offset"`
	PageSize int `json:"page_size"`
}

type FileStatResponse struct {
	BasicResponse

	FileName string    `json:"file_name"`
	PickCode string    `json:"pick_code"`
	Sha1     string    `json:"sha1"`
	IsFile   StringInt `json:"file_category"`

	CreateTime StringInt64 `json:"ptime"`
	UpdateTime StringInt64 `json:"utime"`
	AccessTime int64       `json:"open_time"`

	Paths []*DirInfo `json:"paths"`

	FileCount StringInt `json:"count"`
	DirCount  StringInt `json:"folder_count"`
}
