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
	// Available size in bytes.
	Avail int64

	// Human-readable total size.
	FormatSize string
	// Human-readable used size.
	FormatUsed string
	// Human-readable remain size.
	FormatAvail string
}

// StorageStat gets storage size information.
func (a *Agent) StorageStat(info *StorageInfo) (err error) {
	resp := &webapi.BasicResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiIndexInfo, nil, nil, resp); err != nil {
		return err
	}
	result := webapi.IndexData{}
	if err = resp.Decode(&result); err != nil {
		return
	}
	info.Size = int64(result.Space.Total.Size)
	info.Used = int64(result.Space.Used.Size)
	info.Avail = int64(result.Space.Remain.Size)
	info.FormatSize = result.Space.Total.FormatSize
	info.FormatUsed = result.Space.Used.FormatSize
	info.FormatAvail = result.Space.Remain.FormatSize
	return
}

// StorageFormatInfo describes storage space format usage.
type StorageFormatInfo struct {
	// Total size in bytes.
	Size string
	// Used size in bytes.
	Used string
	// Avail size in bytes.
	Avail string
}

// StorageFormatStat gets storage size information format.
func (a *Agent) StorageFormatStat(info *StorageFormatInfo) (err error) {
	resp := &webapi.BasicResponse{}
	if err = a.pc.CallJsonApi(webapi.ApiIndexInfo, nil, nil, resp); err != nil {
		return err
	}
	result := webapi.IndexData{}
	if err = resp.Decode(&result); err != nil {
		return
	}
	info.Size = result.Space.Total.FormatSize
	info.Used = result.Space.Used.FormatSize
	info.Avail = result.Space.Remain.FormatSize
	return
}
