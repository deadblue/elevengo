package client

import (
	"context"
	"io"
)

// Client is the low-level client which executes ApiSpec.
type Client interface {

	// GetUserAgent returns current "User-Agent" value.
	GetUserAgent() string

	// ExportCookies exports cookies for specific URL.
	ExportCookies(url string) map[string]string

	// Get performs an HTTP GET request.
	Get(url string, headers map[string]string, context context.Context) (body io.ReadCloser, err error)

	// CallApi calls an API.
	CallApi(spec ApiSpec, context context.Context) error
}
