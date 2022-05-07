package oss

import "time"

var GMT *time.Location

func Date() string {
	now := time.Now().In(GMT)
	return now.Format(time.RFC1123)
}

func init() {
	GMT, _ = time.LoadLocation("GMT")
}
