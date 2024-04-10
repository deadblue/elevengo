package elevengo

import (
	"context"

	"github.com/deadblue/elevengo/lowlevel/api"
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
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}
	space := spec.Result.SpaceInfo
	info.Size = int64(space.Total.Size)
	info.Used = int64(space.Used.Size)
	info.Avail = int64(space.Remain.Size)
	info.FormatSize = space.Total.SizeFormat
	info.FormatUsed = space.Used.SizeFormat
	info.FormatAvail = space.Remain.SizeFormat
	return
}
