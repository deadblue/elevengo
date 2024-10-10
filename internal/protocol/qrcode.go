package protocol

import (
	"encoding/json"

	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/errors"
)

//lint:ignore U1000 This type is used in generic.
type QrcodeBaseResp struct {
	State         int    `json:"state"`
	ErrorCode1    int    `json:"code"`
	ErrorCode2    int    `json:"errno"`
	ErrorMessage1 string `json:"message"`
	ErrorMessage2 string `json:"error"`

	Data json.RawMessage `json:"data"`
}

func (r *QrcodeBaseResp) Err() error {
	if r.State != 0 {
		return nil
	}
	return errors.Get(r.ErrorCode1, util.NonEmptyString(
		r.ErrorMessage1, r.ErrorMessage2,
	))
}

func (r *QrcodeBaseResp) Extract(v any) error {
	return json.Unmarshal(r.Data, v)
}
