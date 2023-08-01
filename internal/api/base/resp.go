package base

import "encoding/json"

// BasicResp is the basic response for most JSON/JSONP API.
type BasicResp struct {
	// Response state
	State bool `json:"state"`
	// Error code fields
	ErrorCode  json.Number `json:"errno,omitempty"`
	ErrorCode2 int         `json:"errNo,omitempty"`
	ErrorCode3 int         `json:"errcode,omitempty"`
	ErrorCode4 int         `json:"code,omitempty"`
	// Error message fields
	ErrorMessage  string `json:"error,omitempty"`
	ErrorMessage2 string `json:"message,omitempty"`
}

func (r *BasicResp) Err() error {
	if r.State {
		return nil
	}
	errCode := findNonZero(
		mustInt(r.ErrorCode),
		r.ErrorCode2,
		r.ErrorCode3,
		r.ErrorCode4,
	)
	return GetError(errCode)
}

func mustInt(n json.Number) int {
	if i64, err := n.Int64(); err == nil {
		return int(i64)
	} else {
		return 0
	}
}

func findNonZero(code ...int) int {
	for _, c := range code {
		if c != 0 {
			return c
		}
	}
	return 0
}

type _ApiResp interface {
	Err() error
}

func checkError(r any) error {
	if ar, ok := r.(_ApiResp); ok {
		return ar.Err()
	}
	return nil
}

type GenericResp[D any] struct {
	BasicResp
	Data D `json:"data"`
}
