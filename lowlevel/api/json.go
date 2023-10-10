package api

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/client"
)

type _ApiResp interface {
	// Err returns an error if API calling failed.
	Err() error
}

type _DataApiResp interface {
	// Extract extracts result from response to |v|.
	Extract(v any) error
}

// _JsonApiSpec is the base spec for all JSON ApiSpec.
//
// Type parameters:
//   - D: Result type.
//   - R: Response type.
type _JsonApiSpec[D, R any] struct {
	_BasicApiSpec

	// Request parameters in form
	form util.Params

	// The API result, its value will be filled after |Parse| called.
	Result D
}

func (s *_JsonApiSpec[D, R]) Init(baseUrl string) {
	s._BasicApiSpec.Init(baseUrl)
	s.form = util.Params{}
}

func (s *_JsonApiSpec[D, R]) Payload() client.Payload {
	if len(s.form) == 0 {
		return nil
	} else {
		return client.WwwFormPayload(s.form.Encode())
	}
}

func (s *_JsonApiSpec[D, R]) Parse(r io.Reader) (err error) {
	jd, resp := json.NewDecoder(r), any(new(R))
	if err = jd.Decode(resp); err != nil {
		return
	}
	// Check response error
	ar := resp.(_ApiResp)
	if err = ar.Err(); err != nil {
		return
	}
	// Extract data
	if dr, ok := resp.(_DataApiResp); ok {
		err = dr.Extract(&s.Result)
	}
	return
}

func (s *_JsonApiSpec[D, R]) FormSetAll(params map[string]string) {
	for key, value := range params {
		s.form.Set(key, value)
	}
}

// _JsonpApiSpec is the base spec for all JSON-P ApiSpec.
//
// Type parameters:
//   - D: Result type.
//   - R: Response type.
type _JsonpApiSpec[D, R any] struct {
	_BasicApiSpec

	// The API result, its value will be filled after |Parse| called.
	Result D
}

func (s *_JsonpApiSpec[D, R]) Init(baseUrl, cb string) {
	s._BasicApiSpec.Init(baseUrl)
	s.query.Set("callback", cb)
}

func (s *_JsonpApiSpec[D, R]) Payload() client.Payload {
	return nil
}

func (s *_JsonpApiSpec[D, R]) Parse(r io.Reader) (err error) {
	// Read response
	var body []byte
	if body, err = io.ReadAll(r); err != nil {
		return
	}
	// Find JSON content
	left, right := bytes.IndexByte(body, '('), bytes.LastIndexByte(body, ')')
	if left < 0 || right < 0 {
		return &json.SyntaxError{Offset: 0}
	}
	// Parse JSON
	resp := any(new(R))
	if err = json.Unmarshal(body[left+1:right], resp); err != nil {
		return
	}
	// Force convert resp to ApiResp
	ar := resp.(_ApiResp)
	if err = ar.Err(); err != nil {
		return
	}
	// Extract data
	if dr, ok := resp.(_DataApiResp); ok {
		err = dr.Extract(&s.Result)
	}
	return
}

// _StandardApiSpec is the base spec for all standard JSON API specs.
//
// Type parameters:
//   - D: Result type.
type _StandardApiSpec[D any] struct {
	_JsonApiSpec[D, protocol.StandardResp]
}

// VoidResult describes a void result.
type VoidResult struct{}

// _VoidApiSpec is the base spec for all JSON API specs which has no result.
type _VoidApiSpec struct {
	_JsonApiSpec[VoidResult, protocol.BasicResp]
}
