package internal

type SizeInfo struct {
	Size       float64 `json:"size"`
	SizeFormat string  `json:"size_format"`
}

type FileIndexResult struct {
	_BaseResult
	Data struct {
		SpaceInfo struct {
			AllTotal  SizeInfo `json:"all_total"`
			AllRemain SizeInfo `json:"all_remain"`
			AllUsed   SizeInfo `json:"all_use"`
		} `json:"space_info"`
	} `json:"data"`
}

type Path struct {
	AreaId     IntString `json:"aid"`
	CategoryId IntString `json:"cid"`
	ParentId   IntString `json:"pid"`
	Name       string    `json:"name"`
}

type File struct {
	FileId     string `json:"fid"`
	CategoryId string `json:"cid"`
	ParentId   string `json:"pid"`
	Name       string `json:"n"`
	Size       int64  `json:"s"`
	PickCode   string `json:"pc"`
	Sha1       string `json:"sha"`
	CreateTime string `json:"tp"`
	UpdateTime string `json:"te"`
}

type FileListResult struct {
	_BaseResult
	Count    int    `json:"count"`
	SysCount int    `json:"sys_count"`
	Path     []Path `json:"path"`
	Data     []File `json:"data"`
}
