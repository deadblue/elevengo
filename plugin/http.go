package plugin

import "net/http"

// HttpClient declares the interface which an HTTP client should
// implement, a well-known implementation is `http.Client`, also
// developer can implement it by himself.
type HttpClient interface {

	// Do sends HTTP request and returns HTTP responses.
	Do(req *http.Request) (resp *http.Response, err error)
}

// HttpClientWithJar declares interface for developer, who uses
// self-implemented HttpClient instead of `http.Client`, and
// manages cookie himself.
type HttpClientWithJar interface {
	HttpClient

	// Jar returns client managed cookie jar.
	Jar() http.CookieJar
}
