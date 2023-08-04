package base

import (
	"encoding/json"
	"strconv"
)

// IntNumber uses for JSON field which maybe a string or an integer number.
type IntNumber int64

func (n *IntNumber) UnmarshalJSON(b []byte) (err error) {
	var i int64
	if b[0] == '"' {
		var s string
		if err = json.Unmarshal(b, &s); err == nil {
			i, _ = strconv.ParseInt(s, 10, 64)
		}
	} else {
		err = json.Unmarshal(b, &i)
	}
	if err == nil {
		*n = IntNumber(i)
	}
	return
}

func (n IntNumber) Int64() int64 {
	return int64(n)
}

func (n IntNumber) Int() int {
	return int(n)
}

func (n IntNumber) String() string {
	return strconv.FormatInt(int64(n), 10)
}

// FloatNumner uses for JSON field which maybe a string or an float number.
type FloatNumner float64

func (n *FloatNumner) UnmarshalJSON(b []byte) (err error) {
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
		*n = FloatNumner(f)
	}
	return
}

func (n FloatNumner) Float64() float64 {
	return float64(n)
}

type Boolean bool

func (b *Boolean) UnmarshalJSON(data []byte) (err error) {
	var v bool
	switch data[0] {
	case 'f':
		v = false
	case 't':
		v = true
	case '"':
		var s string
		if err = json.Unmarshal(data, &s); err == nil {
			v = s != ""
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var i int
		if err = json.Unmarshal(data, &i); err == nil {
			v = i != 0
		}
	}
	*b = Boolean(v)
	return
}
