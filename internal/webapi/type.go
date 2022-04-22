package webapi

import (
	"encoding/json"
	"strconv"
)

// StringInt uses for json field which maybe a string or an int.
type StringInt int64

func (v *StringInt) UnmarshalJSON(b []byte) (err error) {
	var i int
	if b[0] == '"' {
		var s string
		if err = json.Unmarshal(b, &s); err == nil {
			i, err = strconv.Atoi(s)
		}
	} else {
		err = json.Unmarshal(b, &i)
	}
	if err == nil {
		*v = StringInt(i)
	}
	return
}

// StringInt64 uses for json field which maybe a string or an int64.
type StringInt64 int64

func (v *StringInt64) UnmarshalJSON(b []byte) (err error) {
	var i int64
	if b[0] == '"' {
		var s string
		if err = json.Unmarshal(b, &s); err == nil {
			i, err = strconv.ParseInt(s, 10, 64)
		}
	} else {
		err = json.Unmarshal(b, &i)
	}
	if err == nil {
		*v = StringInt64(i)
	}
	return
}

// StringFloat64 uses for json field which maybe a string or a float64.
type StringFloat64 float64

func (v *StringFloat64) UnmarshalJSON(b []byte) (err error) {
	var f float64
	if b[0] == '"' {
		var s string
		if err = json.Unmarshal(b, &s); err == nil {
			f, err = strconv.ParseFloat(s, 64)
		}
	} else {
		err = json.Unmarshal(b, &f)
	}
	if err == nil {
		*v = StringFloat64(f)
	}
	return
}

type IntString string

func (v *IntString) UnmarshalJSON(b []byte) (err error) {
	var s string
	if b[0] == '"' {
		err = json.Unmarshal(b, &s)
	} else {
		var i int
		if err = json.Unmarshal(b, &i); err == nil {
			s = strconv.Itoa(i)
		}
	}
	if err == nil {
		*v = IntString(s)
	}
	return
}
