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
	// Extract extracts result from response to |v|.
	Extract(v any) error
}

// JsonApiSpec is the base spec for all JSON ApiSpec.
//
// Type parameters:
//   - D: Result type.
//   - R: Response type.
type JsonApiSpec[D, R any] struct {
	_BaseApiSpec
	form url.Values
	// The API result, its value will be filled after |Parse| called.
	Result D
}

func (s *JsonApiSpec[D, R]) Init(baseUrl string) {
	s._BaseApiSpec.Init(baseUrl)
	s.form = url.Values{}
}

func (s *JsonApiSpec[D, R]) Payload() protocol.Payload {
	if len(s.form) == 0 {
		return nil
	} else {
		return wwwFormPayload(s.form.Encode())
	}
}

func (s *JsonApiSpec[D, R]) Parse(r io.Reader) (err error) {
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
		err = dr.Extract(&s.Result)
	}
	return
}

func (s *JsonApiSpec[D, R]) FormSetAll(params map[string]string) {
	for key, value := range params {
		s.form.Set(key, value)
	}
}

func (s *JsonApiSpec[D, R]) FormSet(key, value string) {
	s.form.Set(key, value)
}

func (s *JsonApiSpec[D, R]) FormSetInt(key string, value int) {
	s.form.Set(key, strconv.Itoa(value))
}

func (s *JsonApiSpec[D, R]) FormSetInt64(key string, value int64) {
	s.form.Set(key, strconv.FormatInt(value, 10))
}

// JsonpApiSpec is the base spec for all JSON-P ApiSpec.
//
// Type parameters:
//   - D: Result type.
//   - R: Response type.
type JsonpApiSpec[D, R any] struct {
	_BaseApiSpec
	// The API result, its value will be filled after |Parse| called.
	Result D
}

func (s *JsonpApiSpec[D, R]) Init(baseUrl, cb string) {
	s._BaseApiSpec.Init(baseUrl)
	s.QuerySet("callback", cb)
}

func (s *JsonpApiSpec[D, R]) Payload() protocol.Payload {
	return nil
}

func (s *JsonpApiSpec[D, R]) Parse(r io.Reader) (err error) {
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
		err = dr.Extract(&s.Result)
	}
	return
}
