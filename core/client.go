package core

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// HTTP client
type HttpClient interface {
	// Get content from a URL
	Get(url string, qs QueryString) ([]byte, error)
	// Call a JSON API with parameters
	JsonApi(url string, qs QueryString, form Form, result interface{}) (err error)
	// Call a JSON-P API with parameters
	JsonpApi(url string, qs QueryString, result interface{}) (err error)
	// Get cookies for URL.
	Cookies(url string) (cookies map[string]string)
	// Set cookies for URL.
	SetCookies(url, domain string, cookies map[string]string)
}

func NewHttpClient(headers http.Header) HttpClient {
	// Make a copy of the default tranport
	tp := http.DefaultTransport.(*http.Transport).Clone()
	// Adjust some parameters
	tp.ForceAttemptHTTP2 = true
	tp.ExpectContinueTimeout = 5 * time.Second
	tp.MaxIdleConnsPerHost = 10
	tp.MaxIdleConns = 50
	tp.IdleConnTimeout = 30 * time.Second
	tp.DialContext = (&net.Dialer{
		Timeout:   0,
		KeepAlive: 30 * time.Second,
	}).DialContext
	// Make http.Client
	jar, _ := cookiejar.New(nil)
	hc := &http.Client{
		Jar:       jar,
		Transport: tp,
		Timeout:   0,
	}
	return &implHttpClient{
		hc:   hc,
		jar:  jar,
		hdrs: headers,
	}
}
