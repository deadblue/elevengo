package plugin

import "net/http"

// HttpClient declares the interface which an HTTP client should
// implement, a well-known implementation is `http.Client`, also
// developer can implement it by himself.
type HttpClient interface {

	// Do sends HTTP request and returns HTTP responses.
	Do(req *http.Request) (resp *http.Response, err error)
}

type HttpClientWithJar interface {
	HttpClient

	// Jar returns cookie jar
	Jar() http.CookieJar
}
