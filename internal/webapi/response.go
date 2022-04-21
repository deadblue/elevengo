package webapi

import (
	"encoding/json"
)

type BasicResponse struct {
	// Response state
	State bool `json:"state"`
	// Error code
	ErrorCode  int `json:"errno,omitempty"`
	ErrorCode2 int `json:"errNo,omitempty"`
	// Error message
	ErrorMessage string `json:"error,omitempty"`
	// Response data
	Data json.RawMessage `json:"data,omitempty"`
}

func (r *BasicResponse) Err() error {
	if !r.State {
		code := r.ErrorCode
		if code == 0 {
			code = r.ErrorCode2
		}
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
