package webapi

var (
	DirOrderModes = []string{
		FileOrderByTime, FileOrderByType, FileOrderBySize, FileOrderByName,
	}
)

type DirMakeResponse struct {
	BasicResponse

	AreaId IntString `json:"aid"`

	CategoryId   string `json:"cid"`
	CategoryName string `json:"cname"`

	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
}

type DirLocateResponse struct {
	BasicResponse
	DirId     IntString `json:"id"`
	IsPrivate StringInt `json:"is_private"`
}
