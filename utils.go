package elevengo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type RequestParameters struct {
	values *url.Values
}

func newRequestParameters() *RequestParameters {
	return &RequestParameters{values: &url.Values{}}
}

func (rp *RequestParameters) With(key, value string) *RequestParameters {
	rp.values.Set(key, value)
	return rp
}

func (rp *RequestParameters) WithStrings(key string, values ...string) *RequestParameters {
	for index, value := range values {
		k := fmt.Sprintf("%s[%d]", key, index)
		rp.values.Set(k, value)
	}
	return rp
}

func (rp *RequestParameters) WithInt(key string, value int) *RequestParameters {
	rp.values.Set(key, strconv.Itoa(value))
	return rp
}

func (rp *RequestParameters) WithInt64(key string, value int64) *RequestParameters {
	rp.values.Set(key, strconv.FormatInt(value, 10))
	return rp
}

func (rp *RequestParameters) QueryString() string {
	return rp.values.Encode()
}

func (rp *RequestParameters) FormData() io.Reader {
	return strings.NewReader(rp.QueryString())
}

type NumberString string

func (ns *NumberString) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		var s string
		err := json.Unmarshal(b, &s)
		if err != nil {
			return err
		} else {
			*ns = NumberString(s)
		}
	} else {
		var n int
		err := json.Unmarshal(b, &n)
		if err != nil {
			return err
		} else {
			*ns = NumberString(strconv.Itoa(n))
		}
	}
	return nil
}
