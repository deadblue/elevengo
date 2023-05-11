package webapi

type OfflineTask struct {
	InfoHash   string  `json:"info_hash"`
	Name       string  `json:"name"`
	Size       int64   `json:"size"`
	FileId     string  `json:"file_id"`
	DirId      string  `json:"wp_path_id"`
	Url        string  `json:"url"`
	Status     int     `json:"status"`
	Percent    float64 `json:"percentDone"`
	AddTime    int64   `json:"add_time"`
	UpdateTime int64   `json:"last_update"`
}

type OfflineBasicResponse struct {
	State     bool `json:"state"`
	ErrorNum  int  `json:"errno"`
	ErrorCode int  `json:"errcode"`
}

func (r *OfflineBasicResponse) Err() error {
	if r.State {
		return nil
	}
	code := r.ErrorNum
	if code == 0 {
		code = r.ErrorCode
	}
	return getError(code)
}

type OfflineListResponse struct {
	OfflineBasicResponse

	Tasks []*OfflineTask `json:"tasks"`

	TaskCount int `json:"count"`
	PageIndex int `json:"page"`
	PageCount int `json:"page_count"`
	PageSize  int `json:"page_row"`

	QuotaTotal  int `json:"total"`
	QuotaRemain int `json:"quota"`
}

type OfflineAddUrlResponse struct {
	OfflineBasicResponse

	InfoHash  string `json:"info_hash"`
	Name      string `json:"name"`
	Url       string `json:"url"`
}

type OfflineAddUrlsResponse struct {
	OfflineBasicResponse

	Result []*OfflineAddUrlResponse `json:"result"`
}
