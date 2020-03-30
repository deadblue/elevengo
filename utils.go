package elevengo

import (
	"encoding/json"
	"strconv"
)

type NumberString string

func (ns *NumberString) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		var s string
		err := json.Unmarshal(b, &s)
		if err != nil {
			return err
		} else {
			*ns = NumberString(s)
		}
	} else {
		var n int
		err := json.Unmarshal(b, &n)
		if err != nil {
			return err
		} else {
			*ns = NumberString(strconv.Itoa(n))
		}
	}
	return nil
}

type SortOption struct {
	flag string
	asc  bool
}

func (so *SortOption) OrderByTime() *SortOption {
	so.flag = orderFlagTime
	return so
}
func (so *SortOption) OrderByName() *SortOption {
	so.flag = orderFlagName
	return so
}
func (so *SortOption) OrderBySize() *SortOption {
	so.flag = orderFlagSize
	return so
}
func (so *SortOption) Asc() *SortOption {
	so.asc = true
	return so
}
func (so *SortOption) Desc() *SortOption {
	so.asc = false
	return so
}
