package webapi

import (
	"encoding/json"
	"fmt"
)

type BasicResponse struct {
	// Response state
	State bool `json:"state"`
	// Error code
	ErrorCode int `json:"errno,omitempty"`
	// Response data
	Data json.RawMessage `json:"data,omitempty"`
}

func (r *BasicResponse) Ok() bool {
	return r.State
}

func (r *BasicResponse) Err() error {
	if !r.State {
		return fmt.Errorf("api error: %d", r.ErrorCode)
	} else {
		return nil
	}
}

func (r *BasicResponse) Decode(result interface{}) error {
	if result != nil {
		return json.Unmarshal(r.Data, result)
	}
	return nil
}
