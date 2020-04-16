package util

import "time"

func ParseUnixTime(s string) *time.Time {
	t := time.Unix(MustParseInt(s), 0)
	return &t
}
