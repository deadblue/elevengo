package webapi

import (
	"encoding/json"
)

type ApiResponse interface {
	Err() error
}

type BasicResponse struct {
	// Response state
	State bool `json:"state"`
	// Error code
	ErrorCode  StringInt `json:"errno,omitempty"`
	ErrorCode2 int       `json:"errNo,omitempty"`
	ErrorCode3 int       `json:"code,omitempty"`
	// Error message
	ErrorMessage  string `json:"error,omitempty"`
	ErrorMessage2 string `json:"message,omitempty"`
	// Response data
	Data json.RawMessage `json:"data,omitempty"`
}

func findNonZero(code ...int) int {
	for _, c := range code {
		if c != 0 {
			return c
		}
	}
	return 0
}

func (r *BasicResponse) Err() error {
	if !r.State {
		code := findNonZero(
			int(r.ErrorCode), r.ErrorCode2, r.ErrorCode3)
		return getError(code)
	}
	return nil
}

func (r *BasicResponse) Decode(result interface{}) error {
	if result != nil {
		return json.Unmarshal(r.Data, result)
	}
	return nil
}
