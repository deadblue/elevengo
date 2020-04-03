package internal

type SpaceSizeInfo struct {
	Size       float64 `json:"size"`
	SizeFormat string  `json:"size_format"`
}

type FileIndexResult struct {
	BaseApiResult
	Error string `json:"error"`
	Data  struct {
		SpaceInfo struct {
			AllTotal  SpaceSizeInfo `json:"all_total"`
			AllRemain SpaceSizeInfo `json:"all_remain"`
			AllUsed   SpaceSizeInfo `json:"all_use"`
		} `json:"space_info"`
	} `json:"data"`
}

type FilePath struct {
	AreaId     IntString `json:"aid"`
	CategoryId IntString `json:"cid"`
	ParentId   IntString `json:"pid"`
	Name       string    `json:"name"`
}

type FileData struct {
	AreaId     string      `json:"aid"`
	FileId     string      `json:"fid"`
	CategoryId string      `json:"cid"`
	ParentId   string      `json:"pid"`
	Name       string      `json:"n"`
	Size       StringInt64 `json:"s"`
	CreateTime string      `json:"tp"`
	UpdateTime string      `json:"te"`
	PickCode   string      `json:"pc"`
	Sha1       string      `json:"sha"`
}

type FileListResult struct {
	BaseApiResult
	Error    string      `json:"error"`
	Count    int         `json:"count"`
	SysCount int         `json:"sys_count"`
	Path     []*FilePath `json:"path"`
	Data     []*FileData `json:"data"`
}

type FileSearchResult struct {
	BaseApiResult
	Error     string      `json:"error"`
	ErrorCode int         `json:"errCode"`
	Count     int         `json:"count"`
	Offset    int         `json:"offset"`
	PageSize  int         `json:"page_size"`
	Data      []*FileData `json:"data"`
}

type FileOperateResult struct {
	BaseApiResult
	ErrorCode IntString `json:"errno"`
	Error     string    `json:"error"`
}

type CategoryAddResult struct {
	FileOperateResult
	AreaId       int    `json:"aid"`
	CategoryId   string `json:"cid"`
	CategoryName string `json:"cname"`
	FileId       string `json:"file_id"`
	FileName     string `json:"file_name"`
}
