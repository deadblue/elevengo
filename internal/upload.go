package internal

type UploadInitResult struct {
	Host        string `json:"host"`
	Policy      string `json:"policy"`
	AccessKeyId string `json:"accessid"`
	ObjectKey   string `json:"object"`
	Callback    string `json:"callback"`
	Signature   string `json:"signature"`
	Expire      int64  `json:"expire"`
}

type UploadData struct {
	AreaId     int    `json:"aid"`
	CategoryId string `json:"cid"`
	FileId     string `json:"file_id"`
	FileName   string `json:"file_name"`
	FileSize   string `json:"file_size"`
	FileTime   int64  `json:"file_ptime"`
	FileStatus int    `json:"file_status"`
	FileType   int    `json:"file_type"`
	IsVideo    int    `json:"is_video"`
	PickCode   string `json:"pick_code"`
	Sha1       string `json:"sha1"`
}

type UploadResult struct {
	BaseApiResult
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *UploadData `json:"data"`
}
