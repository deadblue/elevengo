package util

import "time"

func ParseUnixTime(s string) time.Time {
	return time.Unix(MustParseInt(s), 0)
}
