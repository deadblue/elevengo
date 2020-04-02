package internal

type DownloadInfoResult struct {
	BaseApiResult
	UserId      int    `json:"user_id"`
	PickCode    string `json:"pick_code"`
	FileId      string `json:"file_id"`
	FileName    string `json:"file_name"`
	FileSize    string `json:"file_size"`
	FileUrl     string `json:"file_url"`
	IsVip       int    `json:"is_vip"`
	IsSnap      int    `json:"is_snap"`
	Is115Chrome int    `json:"is_115chrome"`
	MessageCode int    `json:"msg_code"`
	Message     string `json:"msg"`
}
