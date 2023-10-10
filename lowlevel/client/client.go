package client

import "io"

// Client is the low-level client which comm
type Client interface {

	// SetUserAgent set "User-Agent" value which is used in request header.
	SetUserAgent(ua string)

	// GetUserAgent returns current "User-Agent" value.
	GetUserAgent() string

	// ImportCookies imports cookies for specific domains.
	ImportCookies(cookies map[string]string, domains ...string)

	// ExportCookies exports cookies for specific URL.
	ExportCookies(url string) map[string]string

	// Get performs an HTTP GET request.
	Get(url string, headers map[string]string) (body io.ReadCloser, err error)

	// Get performs an HTTP POST request.
	Post(url string, payload Payload) (body io.ReadCloser, err error)

	// CallApi calls an API.
	CallApi(spec ApiSpec) error
}
