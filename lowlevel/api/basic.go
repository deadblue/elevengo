package api

import (
	"strings"

	"github.com/deadblue/elevengo/internal/util"
)

/*
|ApiSpec| inheritance tree:
                    _____________|_BasicApiSpec|___________
                   /                    |                  \
            |_JsonApiSpec|       |_JsonpApiSpec|      |_M115ApiSpec|
		      /        \
|_StandardApiSpec|   |_VoidApiSpec|
*/

// _BasicApiSpec is the base struct for all |client.ApiSpec| implementations.
type _BasicApiSpec struct {
	// Is crypto
	crypto bool

	// Base url
	baseUrl string

	// Query paramsters
	query util.Params
}

// Init initializes _BasicApiSpec, all child structs should call this method
// before use.
func (s *_BasicApiSpec) Init(baseUrl string) {
	s.baseUrl = baseUrl
	s.query = util.Params{}
}

// IsCrypto implements `ApiSpec.IsCrypto`
func (s *_BasicApiSpec) IsCrypto() bool {
	return s.crypto
}

// SetCryptoKey implements `ApiSpec.SetCryptoKey`
func (s *_BasicApiSpec) SetCryptoKey(key string) {
	s.query.Set("k_ec", key)
}

// Url implements `ApiSpec.Url`
func (s *_BasicApiSpec) Url() string {
	if len(s.query) == 0 {
		return s.baseUrl
	} else {
		qs := s.query.Encode()
		if strings.ContainsRune(s.baseUrl, '?') {
			return s.baseUrl + "&" + qs
		} else {
			return s.baseUrl + "?" + qs
		}
	}
}
