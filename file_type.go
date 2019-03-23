package elevengo

type OrderFlag string

const (
	OrderByTime OrderFlag = "user_ptime"
	OrderByName OrderFlag = "file_name"
	OrderBySize OrderFlag = "file_size"
)

const (
	FileListMinLimit = 10
	FileListMaxLimit = 1000
)

type _BasicResult struct {
	State       bool    `json:"state"`
	ErrorNo     int     `json:"errNo"`
	ErrorType   *string `json:"errtype"`
	Error       *string `json:"error"`
	MessageCode int     `json:"msg_code"`
	Message     *string `json:"msg"`
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
