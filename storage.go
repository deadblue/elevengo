package elevengo

import (
	"github.com/deadblue/elevengo/internal/webapi"
)

// StorageInfo describes storage space usage.
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
	resp := &webapi.BasicResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiIndexInfo, nil, nil, resp); err != nil {
		return err
	}
	result := webapi.IndexData{}
	if err = resp.Decode(&result); err != nil {
		return
	}
	info.Size = int64(result.Space.Total.Size)
	info.Used = int64(result.Space.Used.Size)
	info.Avail = int64(result.Space.Remain.Size)
	return
}
