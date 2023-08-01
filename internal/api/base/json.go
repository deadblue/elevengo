package base

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"

	"github.com/deadblue/elevengo/internal/protocol"
)

type JsonApiSpec[R any] struct {
	_BaseApiSpec
	form url.Values
	Resp R
}

func (s *JsonApiSpec[R]) Init(baseUrl string) {
	s._BaseApiSpec.Init(baseUrl)
	s.form = url.Values{}
}

func (s *JsonApiSpec[R]) Payload() protocol.Payload {
	if len(s.form) == 0 {
		return nil
	} else {
		return wwwFormPayload(s.form.Encode())
	}
}

func (s *JsonApiSpec[R]) Parse(r io.Reader) (err error) {
	jd, respPtr := json.NewDecoder(r), &s.Resp
	if err = jd.Decode(respPtr); err != nil {
		return
	}
	return checkError(respPtr)
}

func (s *JsonApiSpec[R]) FormSet(key, value string) {
	s.form.Set(key, value)
}

type JsonpApiSpec[R any] struct {
	_BaseApiSpec
	Resp R
}

func (s *JsonpApiSpec[D]) Init(baseUrl, cb string) {
	s._BaseApiSpec.Init(baseUrl)
	s.QuerySet("callback", cb)
}

func (s *JsonpApiSpec[D]) Payload() protocol.Payload {
	return nil
}

func (s *JsonpApiSpec[D]) Parse(r io.Reader) (err error) {
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
	if err = json.Unmarshal(body[left+1:right], &s.Resp); err != nil {
		return
	}
	return checkError(&s.Resp)
}
