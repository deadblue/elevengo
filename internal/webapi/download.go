package webapi

import "encoding/json"

type DownloadRequest struct {
	Pickcode string `json:"pickcode,omitempty"`
}

type FileDownloadUrl struct {
	Client int    `json:"client"`
	OssId  string `json:"oss_id"`
	Url    string `json:"url"`
}

type DownloadUrl struct {
	Url string
}

func (u *DownloadUrl) UnmarshalJSON(b []byte) (err error) {
	if len(b) > 0 && b[0] == '{' {
		fdl := &FileDownloadUrl{}
		if err = json.Unmarshal(b, fdl); err == nil {
			u.Url = fdl.Url
		}
	}
	return
}

type DownloadInfo struct {
	FileName string      `json:"file_name"`
	FileSize StringInt64 `json:"file_size"`
	PickCode string      `json:"pick_code"`
	Url      DownloadUrl `json:"url"`
}

type DownloadData map[string]*DownloadInfo

func (d DownloadData) IsValid() bool {
	return len(d) > 0
}
