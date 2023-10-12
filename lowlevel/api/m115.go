package api

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/deadblue/elevengo/internal/crypto/m115"
	"github.com/deadblue/elevengo/internal/impl"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/client"
)

type _M115Extractor[D any] func([]byte, *D) error

// _M115ApiSpec is the base API spec for all m115 encoded ApiSpec.
type _M115ApiSpec[D any] struct {
	_BasicApiSpec

	// Cipher key for encrypt/decrypt.
	key m115.Key
	// API parameters.
	params util.Params
	// Custom the extraction process.
	extractor _M115Extractor[D]

	// Final result.
	Result D
}

func (s *_M115ApiSpec[D]) Init(baseUrl string, ex _M115Extractor[D]) {
	s._BasicApiSpec.Init(baseUrl)
	s.extractor = ex
	s.key = m115.GenerateKey()
	s.params = util.Params{}
}

// Payload implements |ApiSpec.Payload|.
func (s *_M115ApiSpec[D]) Payload() client.Payload {
	data, err := json.Marshal(s.params)
	if err != nil {
		return nil
	}
	form := url.Values{}
	form.Set("data", m115.Encode(data, s.key))
	return impl.WwwFormPayload(form.Encode())
}

// Parse implements |ApiSpec.Parse|.
func (s *_M115ApiSpec[D]) Parse(r io.Reader) (err error) {
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
		if s.extractor != nil {
			return s.extractor(result, &s.Result)
		} else {
			return json.Unmarshal(result, &s.Result)
		}

	} else {
		return err
	}
}
