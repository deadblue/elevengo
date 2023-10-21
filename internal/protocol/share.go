package protocol

import (
	"encoding/json"

	"github.com/deadblue/elevengo/lowlevel/errors"
	"github.com/deadblue/elevengo/lowlevel/types"
)

//lint:ignore U1000 This type is used in generic.
type ShareListResp struct {
	BasicResp

	Count int             `json:"count"`
	List  json.RawMessage `json:"list"`
}

func (r *ShareListResp) Extract(v any) (err error) {
	ptr, ok := v.(*types.ShareListResult)
	if !ok {
		return errors.ErrUnsupportedResult
	}
	if err = json.Unmarshal(r.List, &ptr.Items); err != nil {
		return
	}
	ptr.Count = r.Count
	return
}
