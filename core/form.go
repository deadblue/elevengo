package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
)

// Request Form holder
type Form struct {
	isMultiplart bool
	form         url.Values
	partsWriter  *multipart.Writer
	partsBuffer  *bytes.Buffer
}

func NewForm(isMultipart bool) *Form {
	body := Form{
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
func (f *Form) WithString(name, value string) *Form {
	if f.isMultiplart {
		f.partsWriter.WriteField(name, value)
	} else {
		f.form.Set(name, value)
	}
	return f
}
func (f *Form) WithInt(name string, value int) *Form {
	return f.WithString(name, strconv.Itoa(value))
}
func (f *Form) WithInt64(name string, value int64) *Form {
	return f.WithString(name, strconv.FormatInt(value, 10))
}
func (f *Form) WithStrings(name string, value []string) *Form {
	for index, subValue := range value {
		subName := fmt.Sprintf("%s[%d]", name, index)
		f.WithString(subName, subValue)
	}
	return f
}
func (f *Form) WithFile(name, filename string, data io.Reader) *Form {
	if f.isMultiplart {
		w, _ := f.partsWriter.CreateFormFile(name, filename)
		io.Copy(w, data)
	}
	return f
}
func (f *Form) ContentType() string {
	if f.isMultiplart {
		return f.partsWriter.FormDataContentType()
	} else {
		return "application/x-www-form-urlencoded"
	}
}
func (f *Form) Finish() io.Reader {
	if f.isMultiplart {
		f.partsWriter.Close()
		return f.partsBuffer
	} else {
		return strings.NewReader(f.form.Encode())
	}
}
