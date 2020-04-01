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

type UploadResult struct {
}
