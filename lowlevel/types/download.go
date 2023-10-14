package types

import "encoding/json"

type _DownloadUrlProto struct {
	Url    string `json:"url"`
	Client int    `json:"client"`
	OssId  string `json:"oss_id"`
}

type DownloadUrl struct {
	Url string
}

func (u *DownloadUrl) UnmarshalJSON(b []byte) (err error) {
	if len(b) > 0 && b[0] == '{' {
		proto := &_DownloadUrlProto{}
		if err = json.Unmarshal(b, proto); err == nil {
			u.Url = proto.Url
		}
	}
	return
}

type DownloadInfo struct {
	FileName string      `json:"file_name"`
	FileSize json.Number `json:"file_size"`
	PickCode string      `json:"pick_code"`
	Url      DownloadUrl `json:"url"`
}

type DownloadResult map[string]*DownloadInfo
