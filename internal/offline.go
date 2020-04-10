package internal

type OfflineSpaceResult struct {
	BaseApiResult
	Data  float64 `json:"data"`
	Size  string  `json:"size"`
	Limit int64   `json:"limit"`
	Url   string  `json:"url"`
	BtUrl string  `json:"bt_url"`
	Sign  string  `json:"sign"`
	Time  int64   `json:"time"`
}

type OfflineBasicResult struct {
	BaseApiResult
	ErrorCode int    `json:"errcode"`
	ErrorType string `json:"errtype"`
	ErrorMsg  string `json:"error_msg"`
}

type OfflineListResult struct {
	OfflineBasicResult
	Count       int                `json:"count"`
	Page        int                `json:"page"`
	PageCount   int                `json:"page_count"`
	PageSize    int                `json:"page_row"`
	QuotaRemain int                `json:"quota"`
	QuotaTotal  int                `json:"total"`
	Tasks       []*OfflineTaskData `json:"tasks"`
}

type OfflineTaskData struct {
	InfoHash     string  `json:"info_hash"`
	Name         string  `json:"name"`
	Size         int64   `json:"size"`
	Url          string  `json:"url"`
	Status       int     `json:"status"`
	AddTime      int64   `json:"add_time"`
	LeftTime     int64   `json:"left_time"`
	UpdateTime   int64   `json:"last_update"`
	Precent      float64 `json:"percentDone"`
	Move         int     `json:"move"`
	FileId       string  `json:"file_id"`
	DeleteFileId string  `json:"delete_file_id"`
	DeletePath   string  `json:"del_path"`
}

type OfflineAddUrlResult struct {
	OfflineBasicResult
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

type OfflineAddUrlsResult struct {
	OfflineBasicResult
	Result []*OfflineAddUrlResult `json:"result"`
}
