package types

type DirGetIdResult struct {
	FileOperateResult
	Id IntString `json:"id"`
}
