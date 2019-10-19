package elevengo

import (
	"fmt"
	"time"
)

var (
	defaultUserAgent           = fmt.Sprintf("Mozilla/5.0 (Project NTR115; elevengo/%s)", Version)
	defaultDialTimeout         = 10 * time.Second
	defaultIdleTimeout         = 300 * time.Second
	defaultMaxIdleConnsPreHost = 100
	defaultMaxIdleConns        = 200
)

// Options for create `Client`
type Options struct {
	// UserAgent will be sent as "User-Agent" header in each request
	UserAgent string
	// The timeout when connect to API server
	DialTimeout time.Duration
	// The timeout to hold a idle connection before close it
	IdleTimeout time.Duration
	// Max idle connections per host
	MaxIdleConnsPreHost int
	// Max idle connections in total
	MaxIdleConns int
}

func NewOptions() *Options {
	return &Options{
		UserAgent:           defaultUserAgent,
		DialTimeout:         defaultDialTimeout,
		IdleTimeout:         defaultIdleTimeout,
		MaxIdleConnsPreHost: defaultMaxIdleConnsPreHost,
		MaxIdleConns:        defaultMaxIdleConns,
	}
}

func (o *Options) WithUserAgent(value string) *Options {
	o.UserAgent = value
	return o
}

func (o *Options) WithDialTimeout(duration time.Duration) *Options {
	o.DialTimeout = duration
	return o
}

func (o *Options) WithIdleTimeout(duration time.Duration) *Options {
	o.IdleTimeout = duration
	return o
}

func (o *Options) WithMaxIdleConnsPreHost(value int) *Options {
	o.MaxIdleConnsPreHost = value
	return o
}

func (o *Options) WithMaxIdleConns(value int) *Options {
	o.MaxIdleConns = value
	return o
}
