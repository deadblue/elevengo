package core

import (
	"net/http"
	"time"
)

type HttpOpts struct {
	// Cookie jar
	Jar http.CookieJar

	// WaitTimeout
	WaitTimeout time.Duration

	// Handle function before send a request
	BeforeSend func(req *http.Request)

	// Handle function when receive a response
	AfterRecv func(resp *http.Response)
}

// Create a default options
func NewHttpOpts() *HttpOpts {
	return &HttpOpts{
		WaitTimeout: 30 * time.Second,
	}
}
