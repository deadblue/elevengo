package base

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/deadblue/elevengo/internal/crypto/m115"
	"github.com/deadblue/elevengo/internal/protocol"
)

type M115Extractor[D any] func([]byte, *D) error

// M115ApiSpec is the base struct for all m115 ApiSpec.
type M115ApiSpec[D any] struct {
	_BaseApiSpec
	// Cipher key for encrypt/decrypt.
	key m115.Key
	// API parameters.
	params map[string]string
	// Final result.
	Result D
	// Custom the extraction process.
	Extractor M115Extractor[D]
}

func (s *M115ApiSpec[D]) Init(baseUrl string) {
	s._BaseApiSpec.Init(baseUrl)
	s.key = m115.GenerateKey()
	s.params = make(map[string]string)
}

// Payload implements |ApiSpec.Payload|.
func (s *M115ApiSpec[D]) Payload() protocol.Payload {
	data, err := json.Marshal(s.params)
	if err != nil {
		return nil
	}
	form := url.Values{}
	form.Set("data", m115.Encode(data, s.key))
	return wwwFormPayload(form.Encode())
}

// Parse implements |ApiSpec.Parse|.
func (s *M115ApiSpec[D]) Parse(r io.Reader) (err error) {
	jd, resp := json.NewDecoder(r), &StandardResp{}
	if err = jd.Decode(resp); err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return err
	}
	// Decode response data
	var data string
	if err = resp.Extract(&data); err != nil {
		return
	}
	if result, err := m115.Decode(data, s.key); err == nil {
		if s.Extractor != nil {
			return s.Extractor(result, &s.Result)
		} else {
			return json.Unmarshal(result, &s.Result)
		}

	} else {
		return err
	}
}

func (s *M115ApiSpec[D]) ParamSet(key, value string) {
	s.params[key] = value
}

func (s *M115ApiSpec[D]) ParamSetAll(params map[string]string) {
	for key, value := range params {
		s.params[key] = value
	}
}
