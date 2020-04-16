package internal

import (
	"strconv"
	"time"
)

func MustAtoi(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	} else {
		return 0
	}
}

func MustParseInt(s string) int64 {
	if i, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0
	} else {
		return i
	}
}

func MustParseFloat(s string) float64 {
	if f, err := strconv.ParseFloat(s, 64); err != nil {
		return 0
	} else {
		return f
	}
}

func ParseUnixTime(s string) *time.Time {
	t := time.Unix(MustParseInt(s), 0)
	return &t
}
