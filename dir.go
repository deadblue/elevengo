package elevengo

const (
	apiDirGetId = "https://webapi.115.com/files/getid"
)

// DirGetId Retrieves directory ID from full path.
func (a *Agent) DirGetId(path string) (directoryId string, err error) {
	//if strings.HasPrefix(path, "/") {
	//	path = path[1:]
	//}
	//qs := core.NewQueryString().
	//	WithString("path", path)
	//result := &types.DirGetIdResult{}
	//err = a.hc.JsonApi(apiDirGetId, qs, nil, result)
	//if err == nil && result.IsFailed() {
	//	err = types.MakeFileError(int(result.ErrorCode), result.Error)
	//}
	//if err == nil {
	//	if directoryId = string(result.Id); directoryId == "0" {
	//		directoryId, err = "", errDirNotExist
	//	}
	//}
	return
}
