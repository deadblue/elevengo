package elevengo

import "github.com/deadblue/elevengo/internal/types"

// StorageInfo describes storage usage.
type StorageInfo struct {
	// Total size in bytes.
	Size int64
	// Used size in bytes.
	Used int64
	// Avail size in bytes.
	Avail int64
}

// StorageStat gets storage size information.
func (a *Agent) StorageStat(info *StorageInfo) (err error) {
	result := &types.FileIndexResult{}
	err = a.hc.JsonApi(apiFileIndex, nil, nil, result)
	if err == nil && result.IsFailed() {
		err = types.MakeFileError(result.Code, result.Error)
	}
	if err != nil {
		return
	}
	info.Size = int64(result.Data.SpaceInfo.AllTotal.Size)
	info.Used = int64(result.Data.SpaceInfo.AllUsed.Size)
	info.Avail = int64(result.Data.SpaceInfo.AllRemain.Size)
	return
}
