package elevengo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type _QueryString struct {
	values url.Values
}

func newQueryString() *_QueryString {
	return &_QueryString{
		values: url.Values{},
	}
}

func (qs *_QueryString) WithString(name, value string) *_QueryString {
	qs.values.Set(name, value)
	return qs
}

func (qs *_QueryString) WithInt(name string, value int) *_QueryString {
	return qs.WithString(name, strconv.Itoa(value))
}

func (qs *_QueryString) WithTimestamp(name string) *_QueryString {
	value := strconv.FormatInt(time.Now().UnixNano(), 10)
	return qs.WithString(name, value)
}

func (qs *_QueryString) Encode() string {
	return qs.values.Encode()
}

type _Form struct {
	isMultiplart bool
	form         url.Values
	partsWriter  *multipart.Writer
	partsBuffer  *bytes.Buffer
}

func newForm(isMultipart bool) *_Form {
	body := _Form{
		isMultiplart: isMultipart,
	}
	if isMultipart {
		body.partsBuffer = &bytes.Buffer{}
		body.partsWriter = multipart.NewWriter(body.partsBuffer)
	} else {
		body.form = url.Values{}
	}
	return &body
}

func (rb *_Form) WithString(name, value string) *_Form {
	if rb.isMultiplart {
		rb.partsWriter.WriteField(name, value)
	} else {
		rb.form.Set(name, value)
	}
	return rb
}

func (rb *_Form) WithInt(name string, value int) *_Form {
	return rb.WithString(name, strconv.Itoa(value))
}

func (rb *_Form) WithInt64(name string, value int64) *_Form {
	return rb.WithString(name, strconv.FormatInt(value, 10))
}

func (rb *_Form) WithStrings(name string, value []string) *_Form {
	for index, subValue := range value {
		subName := fmt.Sprintf("%s[%d]", name, index)
		rb.WithString(subName, subValue)
	}
	return rb
}

func (rb *_Form) ContentType() string {
	if rb.isMultiplart {
		return rb.partsWriter.FormDataContentType()
	} else {
		return "application/x-www-form-urlencoded"
	}
}

func (rb *_Form) Finish() io.Reader {
	if rb.isMultiplart {
		rb.partsWriter.Close()
		return rb.partsBuffer
	} else {
		return strings.NewReader(rb.form.Encode())
	}
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
