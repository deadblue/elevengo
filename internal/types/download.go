package types

type DownloadRequest struct {
	Pickcode string `json:"pickcode,omitempty"`
}

type DownloadResponse struct {
	State   bool   `json:"state"`
	Error   int    `json:"errno"`
	Message string `json:"msg"`
	Data    string `json:"data"`
}

type DownloadInfo struct {
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
	PickCode string `json:"pick_code"`
	Url      struct {
		Client int    `json:"client"`
		OssId  string `json:"oss_id"`
		Url    string `json:"url"`
	} `json:"url"`
}

type DownloadResult map[string]*DownloadInfo

// DownloadInfoResult holds download information from upstream server.
//
// Deprecated: Old WebAPI result.
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
