package protocol

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

func defaultHttpClient(jar http.CookieJar) *http.Client {
	// Make a copy of the default transport
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// Change some settings
	transport.MaxIdleConnsPerHost = 10
	transport.MaxConnsPerHost = 0
	transport.MaxIdleConns = 100
	transport.IdleConnTimeout = 30 * time.Second
	transport.DialContext = (&net.Dialer{
		Timeout:   0,
		KeepAlive: 30 * time.Second,
	}).DialContext
	if transport.TLSClientConfig != nil {
		transport.TLSClientConfig.InsecureSkipVerify = true
	} else {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	// Make http.Client
	return &http.Client{
		Transport: transport,
		Jar:       jar,
	}
}
