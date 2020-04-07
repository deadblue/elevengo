package core

import (
	"net"
	"net/http"
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
}

type HttpOpts struct {
	// Cookie jar
	Jar http.CookieJar
	// Hook function before send a request
	BeforeSend func(req *http.Request)
	// Hook function when receive a response
	AfterRecv func(resp *http.Response)
}

func NewHttpClient(opts *HttpOpts) HttpClient {
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
	hc := &http.Client{
		Transport: tp,
		Jar:       opts.Jar,
		Timeout:   0,
	}
	return &implHttpClient{
		hc:         hc,
		beforeSend: opts.BeforeSend,
		afterRecv:  opts.AfterRecv,
	}
}
