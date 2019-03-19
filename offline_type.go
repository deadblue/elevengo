package elevengo

type ClearFlag int

const (
	offlineActionTaskList    = "task_lists"
	offlineActionTaskDelete  = "task_del"
	offlineActionTaskClear   = "task_clear"
	offlineActionTaskAddUrl  = "add_task_url"
	offlineActionTaskAddUrls = "add_task_urls"

	ClearComplete ClearFlag = iota
	ClearAll
	ClearFailed
	ClearRunning
	ClearCompleteAndDelete
	ClearAllAndDelete

	ErrorAccountNeedVerify  = 911
	ErrorOfflineIllegalTask = 10003
)

type OfflineSpaceResult struct {
	State bool    `json:"state"`
	Data  float64 `json:"data"`
	Size  string  `json:"size"`
	Url   string  `json:"url"`
	BtUrl string  `json:"bt_url"`
	Limit int64   `json:"limit"`
	Sign  string  `json:"sign"`
	Time  int64   `json:"time"`
}

type OfflineBasicResult struct {
	State        bool    `json:"state"`
	ErrorNo      int     `json:"errno"`
	ErrorCode    int     `json:"errcode"`
	ErrorType    *string `json:"errtype"`
	ErrorMessage *string `json:"error_msg"`
}

type OfflineTask struct {
	InfoHash   string  `json:"info_hash"`
	Status     int     `json:"status"`
	FileId     string  `json:"file_id"`
	Name       string  `json:"name"`
	Size       int64   `json:"size"`
	Percent    float64 `json:"percentDone"`
	AddTime    int64   `json:"add_time"`
	UpdateTime int64   `json:"last_update"`
	Url        string  `json:"url"`
}

type OfflineTaskListResult struct {
	OfflineBasicResult
	Page       int            `json:"page"`
	PageCount  int            `json:"page_count"`
	PageRow    int            `json:"page_row"`
	Count      int            `json:"count"`
	Quota      int            `json:"quota"`
	QuotaTotal int            `json:"total"`
	Tasks      []*OfflineTask `json:"tasks"`
}

type OfflineTaskAddUrlResult struct {
	OfflineBasicResult
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

type OfflineTaskAddUrlsResult struct {
	OfflineBasicResult
	Result []OfflineTaskAddUrlResult `json:"result"`
}
