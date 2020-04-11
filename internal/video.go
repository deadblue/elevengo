package internal

type FileVideoResult struct {
	BaseApiResult
	ErrorCode  int         `json:"errNo"`
	Error      string      `json:"error"`
	FileId     string      `json:"file_id"`
	FileName   string      `json:"file_name"`
	FileSize   StringInt64 `json:"file_size"`
	FileStatus int         `json:"file_status"`
	ParentId   string      `json:"parent_id"`
	Width      StringInt   `json:"width"`
	Height     StringInt   `json:"height"`
	Duration   StringInt64 `json:"play_long"`
	VideoUrl   string      `json:"video_url"`
}
