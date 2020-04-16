package util

import "io"

// Close a io.Closer and ignore its error
func QuietlyClose(c io.Closer) {
	_ = c.Close()
}
