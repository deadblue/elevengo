package protocol

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/errors"
)

// _BasicResp is the basic response for most JSON/JSONP API.
type BasicResp struct {
	// Response state
	State bool `json:"state"`
	// Possible error code fields
	ErrorCode  util.IntNumber `json:"errno,omitempty"`
	ErrorCode2 int            `json:"errNo,omitempty"`
	ErrorCode3 int            `json:"errcode,omitempty"`
	ErrorCode4 int            `json:"errCode,omitempty"`
	ErrorCode5 int            `json:"code,omitempty"`
	// Possible error message fields
	ErrorMessage  string `json:"error,omitempty"`
	ErrorMessage2 string `json:"message,omitempty"`
	ErrorMessage3 string `json:"error_msg,omitempty"`
}

func (r *BasicResp) Err() error {
	if r.State {
		return nil
	}
	return errors.Get(util.NonZero(
		r.ErrorCode.Int(),
		r.ErrorCode2,
		r.ErrorCode3,
		r.ErrorCode4,
		r.ErrorCode5,
	), util.NonEmptyString(
		r.ErrorMessage,
		r.ErrorMessage2,
		r.ErrorMessage3,
	))
}

// StandardResp is the response for all JSON/JSONP APIs with "data" field.
type StandardResp struct {
	BasicResp

	Data json.RawMessage `json:"data"`
}

func (r *StandardResp) Extract(v any) error {
	return json.Unmarshal(r.Data, v)
}
