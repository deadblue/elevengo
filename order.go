package elevengo

import "github.com/deadblue/elevengo/lowlevel/api"

type FileOrder int

const (
	FileOrderByName FileOrder = iota
	FileOrderBySize
	FileOrderByType
	FileOrderByCreateTime
	FileOrderByUpdateTime
	FileOrderByOpenTime
)

var (
	fileOrderNames = []string{
		api.FileOrderByName,
		api.FileOrderBySize,
		api.FileOrderByType,
		api.FileOrderByCreateTime,
		api.FileOrderByUpdateTime,
		api.FileOrderByOpenTime,
	}
	fileOrderCount = len(fileOrderNames)
)

func getOrderName(order FileOrder) string {
	if order < 0 || int(order) >= fileOrderCount {
		return api.FileOrderDefault
	}
	return fileOrderNames[order]
}
