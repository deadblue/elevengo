package internal

import (
	"strconv"
	"time"
)

func MustParseInt(s string) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0
	} else {
		return i
	}
}

func ParseUnixTime(s string) time.Time {
	ts := MustParseInt(s)
	return time.Unix(ts, 0)
}
