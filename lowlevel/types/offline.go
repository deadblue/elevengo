package types

type OfflineTask struct {
	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Url      string `json:"url"`
	AddTime  int64  `json:"add_time"`

	Status     int     `json:"status"`
	Percent    float64 `json:"percentDone"`
	UpdateTime int64   `json:"last_update"`

	FileId string `json:"file_id"`
	DirId  string `json:"wp_path_id"`
}

type OfflineListResult struct {
	PageIndex int
	PageCount int
	PageSize  int

	TaskCount int
	Tasks     []*OfflineTask
}

type OfflineAddResult struct {
	State   bool   `json:"state"`
	ErrNum  int    `json:"errno"`
	ErrCode int    `json:"errcode"`
	ErrType string `json:"errtype"`

	InfoHash string `json:"info_hash"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

type OfflineAddUrlsData struct {
	State   bool `json:"state"`
	ErrNum  int  `json:"errno"`
	ErrCode int  `json:"errcode"`

	Result []*OfflineAddResult `json:"result"`
}

type OfflineAddUrlsResult []*OfflineTask
