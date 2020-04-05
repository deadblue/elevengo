package internal

import "encoding/json"

type StorageSizeInfo struct {
	Size       float64 `json:"size"`
	SizeFormat string  `json:"size_format"`
}

type FileIndexResult struct {
	BaseApiResult
	Error string `json:"error"`
	Data  struct {
		SpaceInfo struct {
			AllTotal  StorageSizeInfo `json:"all_total"`
			AllRemain StorageSizeInfo `json:"all_remain"`
			AllUsed   StorageSizeInfo `json:"all_use"`
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
	AreaId     IntString   `json:"aid"`
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
	Error     string      `json:"error"`
	ErrorCode int         `json:"errNo"`
	Count     int         `json:"count"`
	SysCount  int         `json:"sys_count"`
	Path      []*FilePath `json:"path"`
	Data      []*FileData `json:"data"`
	Order     string      `json:"order"`
	IsAsc     int         `json:"is_asc"`
	Offset    int         `json:"offset"`
	Limit     int         `json:"limit"`
	PageSize  int         `json:"page_size"`
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
	ErrorCode StringInt `json:"errno"`
	Error     string    `json:"error"`
}

type FileAddResult struct {
	FileOperateResult
	AreaId       int    `json:"aid"`
	CategoryId   string `json:"cid"`
	CategoryName string `json:"cname"`
	FileId       string `json:"file_id"`
	FileName     string `json:"file_name"`
}

type FileStatPath struct {
	FileId   IntString `json:"file_id"`
	FileName string    `json:"file_name"`
}

type FileStatData struct {
	Name        string          `json:"file_name"`
	IsFile      string          `json:"file_category"`
	CreateTime  string          `json:"ptime"`
	UpdateTime  string          `json:"utime"`
	FileCount   StringInt       `json:"count"`
	FolderCount StringInt       `json:"folder_count"`
	FormatSize  string          `json:"size"`
	PickCode    string          `json:"pick_code"`
	Sha1        string          `json:"sha1"`
	Paths       []*FileStatPath `json:"paths"`
}

type FileStatResult struct {
	BaseApiResult
	Data *FileStatData
}

func (r *FileStatResult) UnmarshalJSON(data []byte) (err error) {
	if data[0] == '[' {
		*r = FileStatResult{
			BaseApiResult: BaseApiResult{
				State: false,
			},
			Data: nil,
		}
	} else {
		d := &FileStatData{}
		if err = json.Unmarshal(data, d); err == nil {
			*r = FileStatResult{
				BaseApiResult: BaseApiResult{
					State: true,
				},
				Data: d,
			}
		}
	}
	return
}
