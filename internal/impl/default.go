package impl

import (
	"net"
	"net/http"
	"time"
)

const (
	defaultUserAgent = "Mozilla/5.0"
)

func defaultHttpClient(jar http.CookieJar) *http.Client {
	// Make a copy of the default transport
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// Change some settings
	transport.MaxIdleConnsPerHost = 10
	transport.MaxConnsPerHost = 0
	transport.MaxIdleConns = 100
	transport.IdleConnTimeout = 60 * time.Second
	// Setup timeout
	transport.DialContext = (&net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
	}).DialContext
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.ResponseHeaderTimeout = 30 * time.Second
	// Make http.Client
	return &http.Client{
		Transport: transport,
		Jar:       jar,
	}
}
