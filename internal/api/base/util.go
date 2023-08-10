package base

import "encoding/json"

func MustInt(n json.Number) int {
	if i64, err := n.Int64(); err == nil {
		return int(i64)
	} else {
		return 0
	}
}
