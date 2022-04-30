package webapi

type DirGetIdResponse struct {
	BasicResponse
	DirId     IntString `json:"id"`
	IsPrivate StringInt `json:"is_private"`
}
