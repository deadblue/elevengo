package webapi

import (
	"strconv"
	"strings"
	"time"
)

func ParseFileTime(t string) time.Time {

	if isTimestamp(t) {
		i, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return time.Now()
		}

		return time.Unix(i, 0)
	}
	tm, _ := time.Parse("2006-01-02 15:04", t)
	return tm
}

func isTimestamp(str string) bool {
	// 2023-01-10 13:43
	if strings.Contains(str, "-") {
		return false
	}
	return true
}
