package elevengo

// structs for client

type Options struct {
	// The UserAgent
	UserAgent string

	// Do not use proxy in environment
	DisableProxy bool

	// Max idle connections number per host
	MaxIdleConns int

	// Enable debug mode
	// When enabled, the client:
	//  * Does not verify server certificate
	Debug bool
}

type _UserInfo struct {
	UserId string
}

type _OfflineToken struct {
	Sign string
	Time int64
}

// structs for File API

type _BasicResult struct {
	State       bool    `json:"state"`
	ErrorNo     int     `json:"errNo"`
	ErrorType   *string `json:"errtype"`
	Error       *string `json:"error"`
	MessageCode int     `json:"msg_code"`
	Message     *string `json:"msg"`
}

type SortOptions struct {
	Flag *OrderFlag
	Asc  bool
}

type FolderData struct {
	AreaId     NumberString `json:"aid"`
	CategoryId NumberString `json:"cid"`
	ParentId   NumberString `json:"pid"`
	Name       string       `json:"name"`
}

type FileData struct {
	FileId     *string `json:"fid"`
	CategoryId *string `json:"cid"`
	ParentId   *string `json:"pid"`
	Name       string  `json:"n"`
	Size       int     `json:"s"`
	PickCode   string  `json:"pc"`
	Sha1       string  `json:"sha"`
}

type FileListResult struct {
	_BasicResult
	TotalCount int           `json:"count"`
	SysCount   int           `json:"sys_count"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	PageSize   int           `json:"page_size"`
	Path       []*FolderData `json:"path"`
	Data       []*FileData   `json:"data"`
}

type FileSearchResult struct {
	_BasicResult
	TotalCount int         `json:"count"`
	Offset     int         `json:"offset"`
	PageSize   int         `json:"page_size"`
	Folder     *FolderData `json:"folder"`
	Data       []*FileData `json:"data"`
}

type FileInfoData struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	PickCode string `json:"pick_code"`
	Sha1     string `json:"sha1"`
}

type FileInfoResult struct {
	_BasicResult
	Data []*FileInfoData `json:"data"`
}

type FileAddResult struct {
	_BasicResult
	AreaId       NumberString `json:"aid"`
	CategoryId   NumberString `json:"cid"`
	CategoryName string       `json:"cname"`
	FileId       string       `json:"file_id"`
	FileName     string       `json:"file_name"`
}

type CategoryGetResult struct {
	CategoryName string `json:"file_name"`
	FileCount    string `json:"count"`
	FolderCount  string `json:"folder_count"`
	Size         string `json:"size"`
}

type FileDownloadResult struct {
	_BasicResult
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
	Pickcode string `json:"pickcode"`
	FileUrl  string `json:"file_url"`
}

// structs for Offline API

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
	RealFileId string  `json:"delete_file_id"`
	Name       string  `json:"name"`
	Size       int64   `json:"size"`
	Percent    float64 `json:"percentDone"`
	AddTime    int64   `json:"add_time"`
	UpdateTime int64   `json:"last_update"`
	Url        string  `json:"url"`
}

type OfflineListResult struct {
	OfflineBasicResult
	Page       int            `json:"page"`
	PageCount  int            `json:"page_count"`
	PageRow    int            `json:"page_row"`
	Count      int            `json:"count"`
	Quota      int            `json:"quota"`
	QuotaTotal int            `json:"total"`
	Tasks      []*OfflineTask `json:"tasks"`
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

// structs for Captcha API

type CaptchaSession struct {
	Callback  string
	CodeValue []byte
	CodeKeys  []byte
}

type _CaptchaSignResult struct {
	State bool   `json:"state"`
	Sign  string `json:"sign"`
}

// structs for download API

type DownloadCookie struct {
	Name  string
	Value string
}

type DownloadInfo struct {
	Url       string
	UserAgent string
	Cookies   []*DownloadCookie
}

// structs for upload API

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
	State   bool   `json:"state"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *UploadedFile
}

type UploadedFile struct {
	CategoryId string `json:"cid"`
	FileId     string `json:"file_id"`
	FileName   string `json:"file_name"`
	FizeSize   string `json:"file_size"`
	PickCode   string `json:"pick_code"`
	Sha1       string `json:"sha1"`
}
