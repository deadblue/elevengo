package elevengo

type FileData struct {
	FileId      *string `json:"fid"`
	ContainerId *string `json:"cid"`
	ParentId    string  `json:"pid"`
	Name        string  `json:"n"`
	Size        int     `json:"s"`
	PickCode    string  `json:"pc"`
	Sha1        string  `json:"sha"`
}

type FilePath struct {
	Aid  NumberString `json:"aid"`
	Cid  NumberString `json:"cid"`
	Pid  NumberString `json:"pid"`
	Name string       `json:"name"`
}

type BasicFileResult struct {
	State bool `json:"state"`
}

type FileListResult struct {
	BasicFileResult
	ErrorCode int        `json:"errno"`
	Error     string     `json:"error"`
	Count     int        `json:"count"`
	SysCount  int        `json:"sys_count"`
	Offset    int        `json:"offset"`
	Limit     int        `json:"limit"`
	Path      []FilePath `json:"path"`
	Data      []FileData `json:"data"`
}

type FileDownloadResult struct {
	BasicFileResult
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileSize string `json:"file_size"`
	Pickcode string `json:"pickcode"`
	FileUrl  string `json:"file_url"`
}
