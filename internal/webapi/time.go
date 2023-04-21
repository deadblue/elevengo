package webapi

import (
	"strconv"
	"time"
)

const (
	timeLayout = "2006-01-02 15:04"
)

var (
	timeLocation = time.FixedZone("CST", 8*60*60)
)

func ParseFileTime(str string) time.Time {
	if isTimestamp(str) {
		sec, _ := strconv.ParseInt(str, 10, 64)
		return time.Unix(sec, 0)
	} else {
		t, _ := time.ParseInLocation(timeLayout, str, timeLocation)
		return t.UTC()
	}
}

func isTimestamp(str string) bool {
	for _, ch := range str {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
