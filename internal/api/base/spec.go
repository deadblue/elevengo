package base

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

// _BaseApiSpec is the base struct for all protocol.ApiSpec implementations.
// It does not expose to the outside of this package, child struct should use
// |JsonSpec|, |JsonpSpec| or |M115Spec|.
type _BaseApiSpec struct {
	// Is crypto
	crypto bool
	// Base url
	baseUrl string
	// Query paramsters
	query url.Values
}

// Init initializes BasicApiSpec, all child structs should call this method
// before use.
func (s *_BaseApiSpec) Init(baseUrl string) {
	s.baseUrl = baseUrl
	s.query = url.Values{}
}

func (s *_BaseApiSpec) IsCrypto() bool {
	return s.crypto
}

// Url implements ApiSpec.Url()
func (s *_BaseApiSpec) Url() string {
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

func (s *_BaseApiSpec) EnableCrypto() {
	s.crypto = true
}

func (s *_BaseApiSpec) QuerySet(key, value string) {
	s.query.Set(key, value)
}

func (s *_BaseApiSpec) QuerySetInt(key string, value int) {
	s.query.Set(key, strconv.Itoa(value))
}

func (s *_BaseApiSpec) QuerySetInt64(key string, value int64) {
	s.query.Set(key, strconv.FormatInt(value, 10))
}

func (s *_BaseApiSpec) QuerySetNow(key string) {
	now := time.Now().Unix()
	s.query.Set(key, strconv.FormatInt(now, 10))
}

func (s *_BaseApiSpec) QueryGet(key string) string {
	return s.query.Get(key)
}

func (s *_BaseApiSpec) SetBaseUrl(baseUrl string) {
	s.baseUrl = baseUrl
}
