package internal

import (
	"encoding/json"
	"strconv"
)

type IntString string

func (is *IntString) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		var s string
		err := json.Unmarshal(b, &s)
		if err != nil {
			return err
		} else {
			*is = IntString(s)
		}
	} else {
		var n int
		err := json.Unmarshal(b, &n)
		if err != nil {
			return err
		} else {
			*is = IntString(strconv.Itoa(n))
		}
	}
	return nil
}
