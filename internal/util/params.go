package util

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Params map[string]string

func (p Params) Set(key, value string) Params {
	p[key] = value
	return p
}

func (p Params) SetInt(key string, value int) Params {
	p.Set(key, strconv.Itoa(value))
	return p
}

func (p Params) SetInt64(key string, value int64) Params {
	p.Set(key, strconv.FormatInt(value, 10))
	return p
}

func (p Params) SetNow(key string) Params {
	now := time.Now().Unix()
	p.Set(key, strconv.FormatInt(now, 10))
	return p
}

func (p Params) SetAll(params map[string]string) Params {
	for key, value := range params {
		p.Set(key, value)
	}
	return p
}

func (p Params) Encode() string {
	sb := &strings.Builder{}
	isFirst := true
	for key, value := range p {
		if !isFirst {
			sb.WriteRune('&')
		}
		sb.WriteString(url.QueryEscape(key))
		sb.WriteRune('=')
		sb.WriteString(url.QueryEscape(value))
		isFirst = false
	}
	return sb.String()
}
