package elevengo

import (
	"github.com/deadblue/elevengo/internal/api"
)

// StorageInfo describes storage space usage.
type StorageInfo struct {
	// Total size in bytes.
	Size int64
	// Human-readable total size.
	FormatSize string

	// Used size in bytes.
	Used int64
	// Human-readable used size.
	FormatUsed string

	// Available size in bytes.
	Avail int64
	// Human-readable remain size.
	FormatAvail string
}

// StorageStat gets storage size information.
func (a *Agent) StorageStat(info *StorageInfo) (err error) {
	spec := (&api.IndexInfoSpec{}).Init()
	if err = a.pc.ExecuteApi(spec); err != nil {
		return
	}
	result := spec.Data.SpaceInfo
	info.Size = int64(result.Total.Size)
	info.Used = int64(result.Used.Size)
	info.Avail = int64(result.Remain.Size)
	info.FormatSize = result.Total.SizeFormat
	info.FormatUsed = result.Used.SizeFormat
	info.FormatAvail = result.Remain.SizeFormat
	return
}
