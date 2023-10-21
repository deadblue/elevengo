package types

type ImageGetResult struct {
	FileName string `json:"file_name"`
	FileSha1 string `json:"file_sha1"`
	Pickcode string `json:"pick_code"`

	SourceUrl string   `json:"source_url"`
	OriginUrl string   `json:"origin_url"`
	ViewUrl   string   `json:"url"`
	ThumbUrls []string `json:"all_url"`
}
