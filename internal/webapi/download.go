package webapi

type DownloadRequest struct {
	Pickcode string `json:"pickcode,omitempty"`
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
