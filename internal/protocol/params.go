package protocol

import (
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Params holds parameter pairs.
type Params map[string]string

func (p Params) With(name, value string) Params {
	p[name] = value
	return p
}

func (p Params) WithInt(name string, value int) Params {
	return p.With(name, strconv.Itoa(value))
}

func (p Params) WithInt64(name string, value int64) Params {
	return p.With(name, strconv.FormatInt(value, 10))
}

func (p Params) WithNow(name string) Params {
	return p.WithInt64(name, time.Now().Unix())
}

// Encode encodes params into query-string format.
func (p Params) Encode() string {
	buf, isFirst := strings.Builder{}, true
	for name, value := range p {
		if !isFirst {
			buf.WriteRune('&')
		}
		buf.WriteString(url.QueryEscape(name))
		buf.WriteRune('=')
		buf.WriteString(url.QueryEscape(value))
		isFirst = false
	}
	return buf.String()
}

func (p Params) Reader() io.Reader {
	return strings.NewReader(p.Encode())
}
