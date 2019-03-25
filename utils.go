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

func (f *_Form) WithString(name, value string) *_Form {
	if f.isMultiplart {
		f.partsWriter.WriteField(name, value)
	} else {
		f.form.Set(name, value)
	}
	return f
}

func (f *_Form) WithInt(name string, value int) *_Form {
	return f.WithString(name, strconv.Itoa(value))
}

func (f *_Form) WithInt64(name string, value int64) *_Form {
	return f.WithString(name, strconv.FormatInt(value, 10))
}

func (f *_Form) WithStrings(name string, value []string) *_Form {
	for index, subValue := range value {
		subName := fmt.Sprintf("%s[%d]", name, index)
		f.WithString(subName, subValue)
	}
	return f
}

func (f *_Form) WithFile(name, filename string, data io.Reader) *_Form {
	if f.isMultiplart {
		w, _ := f.partsWriter.CreateFormFile(name, filename)
		io.Copy(w, data)
	}
	return f
}

func (f *_Form) ContentType() string {
	if f.isMultiplart {
		return f.partsWriter.FormDataContentType()
	} else {
		return "application/x-www-form-urlencoded"
	}
}

func (f *_Form) Finish() io.Reader {
	if f.isMultiplart {
		f.partsWriter.Close()
		return f.partsBuffer
	} else {
		return strings.NewReader(f.form.Encode())
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
