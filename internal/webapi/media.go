package webapi

type VideoResponse struct {
	BasicResponse
	// Video information
	Width    StringInt     `json:"width"`
	Height   StringInt     `json:"height"`
	Duration StringFloat64 `json:"play_long"`
	VideoUrl string        `json:"video_url"`
	// File information
	FileId     string      `json:"file_id"`
	FileName   string      `json:"file_name"`
	FileSize   StringInt64 `json:"file_size"`
	FileStatus int         `json:"file_status"`
	PickCode   string      `json:"pick_code"`
	Sha1       string      `json:"sha1"`
}
