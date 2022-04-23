package web

import (
	"fmt"
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

func (p Params) WithArray(name string, values []string) Params {
	for i, value := range values {
		key := fmt.Sprintf("%s[%d]", name, i)
		p.With(key, value)
	}
	return p
}

func (p Params) WithMap(name string, value map[string]string) Params {
	for ik, iv := range value {
		ik = fmt.Sprintf("%s[%s]", name, ik)
		p.With(iv, iv)
	}
	return p
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
