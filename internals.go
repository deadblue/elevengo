package elevengo

type _UserInfo struct {
	UserId string
}

type _OfflineToken struct {
	Sign string
	Time int64
}

type _BasicResult struct {
	State       bool    `json:"state"`
	ErrorNo     int     `json:"errNo"`
	ErrorType   *string `json:"errtype"`
	Error       *string `json:"error"`
	MessageCode int     `json:"msg_code"`
	Message     *string `json:"msg"`
}

type _FileDownloadResult struct {
	_BasicResult
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
	Pickcode string `json:"pickcode"`
	FileUrl  string `json:"file_url"`
}

type _UploadInitResult struct {
	AccessKeyId string `json:"accessid"`
	Callback    string `json:"callback"`
	Expire      int    `json:"expire"`
	UploadUrl   string `json:"host"`
	ObjectKey   string `json:"object"`
	Policy      string `json:"policy"`
	Signature   string `json:"signature"`
}

type _UploadResult struct {
	State   bool          `json:"state"`
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *UploadedFile `json:"data"`
}

type _OfflineBasicResult struct {
	State        bool    `json:"state"`
	ErrorNo      int     `json:"errno"`
	ErrorCode    int     `json:"errcode"`
	ErrorType    string  `json:"errtype"`
	ErrorMessage *string `json:"error_msg"`
}

type _OfflineSpaceResult struct {
	State bool    `json:"state"`
	Data  float64 `json:"data"`
	Size  string  `json:"size"`
	Url   string  `json:"url"`
	BtUrl string  `json:"bt_url"`
	Limit int64   `json:"limit"`
	Sign  string  `json:"sign"`
	Time  int64   `json:"time"`
}

type _OfflineGetDirResult struct {
	CategoryId string `json:"cid"`
}

type _TorrentFile struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type _OfflineTorrentResult struct {
	_OfflineBasicResult
	TorrentName string          `json:"torrent_name"`
	InfoHash    string          `json:"info_hash"`
	FileSize    int64           `json:"file_size"`
	FileCount   int             `json:"file_count"`
	FileList    []*_TorrentFile `json:"torrent_filelist_web"`
}

type _CaptchaSignResult struct {
	State bool   `json:"state"`
	Sign  string `json:"sign"`
}
