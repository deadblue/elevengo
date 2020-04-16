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

type Form interface {
	// Add a string value
	WithString(name, value string) Form
	// Add an int value
	WithInt(name string, value int) Form
	// Add an int64 value
	WithInt64(name string, value int64) Form
	// Add a string slice value
	WithStrings(name string, value []string) Form
	// Add a string map value
	WithStringMap(name string, value map[string]string) Form
	// Add a file
	WithFile(name, filename string, stream io.Reader) Form
	// Return content-type of the form
	ContentType() string
	// Archive the form and return the data stream
	// You can not add values after call this method.
	Archive() io.Reader
}

type implForm struct {
	isMP bool
	// for url-encoded form
	v url.Values
	// for multipart form
	w *multipart.Writer
	b *bytes.Buffer
}

func (f *implForm) WithString(name, value string) Form {
	if f.isMP {
		f.w.WriteField(name, value)
	} else {
		f.v.Set(name, value)
	}
	return f
}
func (f *implForm) WithInt(name string, value int) Form {
	return f.WithString(name, strconv.Itoa(value))
}
func (f *implForm) WithInt64(name string, value int64) Form {
	return f.WithString(name, strconv.FormatInt(value, 10))
}
func (f *implForm) WithStrings(name string, value []string) Form {
	for index, subValue := range value {
		subName := fmt.Sprintf("%s[%d]", name, index)
		f.WithString(subName, subValue)
	}
	return f
}
func (f *implForm) WithStringMap(name string, value map[string]string) Form {
	for mapKey, mapVal := range value {
		subName := fmt.Sprintf("%s[%s]", name, mapKey)
		f.WithString(subName, mapVal)
	}
	return f
}
func (f *implForm) WithFile(name, filename string, stream io.Reader) Form {
	if f.isMP {
		w, _ := f.w.CreateFormFile(name, filename)
		io.Copy(w, stream)
	}
	return f
}
func (f *implForm) ContentType() string {
	if f.isMP {
		return f.w.FormDataContentType()
	} else {
		return "application/x-www-form-urlencoded"
	}
}
func (f *implForm) Archive() io.Reader {
	if f.isMP {
		f.w.Close()
		return f.b
	} else {
		return strings.NewReader(f.v.Encode())
	}
}

func NewForm() Form {
	return &implForm{
		isMP: false,
		v:    url.Values{},
	}
}

func NewMultipartForm() Form {
	buf := &bytes.Buffer{}
	return &implForm{
		isMP: true,
		b:    buf,
		w:    multipart.NewWriter(buf),
	}
}
