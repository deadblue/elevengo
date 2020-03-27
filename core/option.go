package core

import "net/http"

type HttpClientOpts struct {
	// Cookie jar
	Jar http.CookieJar

	// Handle function before send a request
	BeforeSend func(req *http.Request)

	// Handle function when receive a response
	AfterRecv func(resp *http.Response)
}
