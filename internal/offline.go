package internal

type OfflineSpaceResult struct {
	State bool    `json:"state"`
	Data  float64 `json:"data"`
	Size  string  `json:"size"`
	Limit int64   `json:"limit"`
	Url   string  `json:"url"`
	BtUrl string  `json:"bt_url"`
	Sign  string  `json:"sign"`
	Time  int64   `json:"time"`
}

type OfflineBasicResult struct {
	State     bool   `json:"state"`
	ErrorCode int    `json:"errcode"`
	ErrorType string `json:"errtype"`
	ErrorMsg  string `json:"error_msg"`
}

type OfflineListResult struct {
	OfflineBasicResult
	Count     int `json:"count"`
	Page      int `json:"page"`
	PageCount int `json:"page_count"`
	PageSize  int `json:"page_row"`
	Tasks     []*OfflineTask
}

type OfflineTask struct {
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	AddTime  int64  `json:"add_time"`
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
