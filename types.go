package elevengo

// Options for create `Client`
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

type Credentials struct {
	UID  string
	CID  string
	SEID string
}

// File info struct
type CloudFile struct {
	ParentId   string
	CategoryId string
	FileId     string
	IsCategory bool
	IsSystem   bool
	Name       string
	Size       int64
	PickCode   string
	Sha1       string
	CreateTime int64
	UpdateTime int64
}

type CategoryInfoResult struct {
	CategoryName string `json:"file_name"`
	Size         string `json:"size"`
	FileCount    string `json:"count"`
	FolderCount  string `json:"folder_count"`
	CreateTime   int64  `json:"ptime"`
	UpdateTime   int64  `json:"utime"`
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

type CaptchaSession struct {
	Callback  string
	CodeValue []byte
	CodeKeys  []byte
}

type DownloadCookie struct {
	Name  string
	Value string
}

type DownloadInfo struct {
	Url       string
	UserAgent string
	Cookies   []*DownloadCookie
}
