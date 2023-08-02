package base

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strconv"

	"github.com/deadblue/elevengo/internal/protocol"
)

type _ApiResp interface {
	// Err returns an error if API failed.
	Err() error
}

type _DataResp interface {
	// Extract extracts information to |v| from response.
	Extract(v any) error
}

// JsonApiSpec is the base spec for all JSON ApiSpec, child implementation
// should provide type parameters R (response type) and D (data type).
type JsonApiSpec[R any, D any] struct {
	_BaseApiSpec
	form url.Values
	Data D
}

func (s *JsonApiSpec[R, D]) Init(baseUrl string) {
	s._BaseApiSpec.Init(baseUrl)
	s.form = url.Values{}
}

func (s *JsonApiSpec[R, D]) Payload() protocol.Payload {
	if len(s.form) == 0 {
		return nil
	} else {
		return wwwFormPayload(s.form.Encode())
	}
}

func (s *JsonApiSpec[R, D]) Parse(r io.Reader) (err error) {
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
	if dr, ok := resp.(_DataResp); ok {
		err = dr.Extract(&s.Data)
	}
	return
}

func (s *JsonApiSpec[R, D]) FormSet(key, value string) {
	s.form.Set(key, value)
}

func (s *JsonApiSpec[R, D]) FormSetInt(key string, value int) {
	s.form.Set(key, strconv.Itoa(value))
}

// JsonpApiSpec is the base spec for all JSON-P ApiSpec, child implementation
// should provide type parameters R (response type) and D (data type).
type JsonpApiSpec[R any, D any] struct {
	_BaseApiSpec
	// Data holds the final result
	Data D
}

func (s *JsonpApiSpec[R, D]) Init(baseUrl, cb string) {
	s._BaseApiSpec.Init(baseUrl)
	s.QuerySet("callback", cb)
}

func (s *JsonpApiSpec[R, D]) Payload() protocol.Payload {
	return nil
}

func (s *JsonpApiSpec[R, D]) Parse(r io.Reader) (err error) {
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
	if dr, ok := resp.(_DataResp); ok {
		err = dr.Extract(&s.Data)
	}
	return
}
