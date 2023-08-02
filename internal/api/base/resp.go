package base

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/api/errors"
)

// BasicResp is the basic response for most JSON/JSONP API.
type BasicResp struct {
	// Response state
	State bool `json:"state"`
	// Possible error code fields
	ErrorCode  json.Number `json:"errno,omitempty"`
	ErrorCode2 int         `json:"errNo,omitempty"`
	ErrorCode3 int         `json:"errcode,omitempty"`
	ErrorCode4 int         `json:"code,omitempty"`
	// Possible error message fields
	ErrorMessage  string `json:"error,omitempty"`
	ErrorMessage2 string `json:"message,omitempty"`
	ErrorMessage3 string `json:"error_msg,omitempty"`
}

func (r *BasicResp) Err() error {
	if r.State {
		return nil
	}
	errCode := findNonZero(
		MustInt(r.ErrorCode),
		r.ErrorCode2,
		r.ErrorCode3,
		r.ErrorCode4,
	)
	return errors.Get(errCode)
}

func findNonZero(code ...int) int {
	for _, c := range code {
		if c != 0 {
			return c
		}
	}
	return 0
}

// func checkError(r any) error {
// 	if ar, ok := r.(_ApiResp); ok {
// 		return ar.Err()
// 	}
// 	return nil
// }

type StandardResp struct {
	BasicResp
	Data json.RawMessage `json:"data"`
}

func (r *StandardResp) Extract(v any) error {
	return json.Unmarshal(r.Data, v)
}