package util

import (
	"mime"
	"path"
)

const (
	DefaultMimeType = "application/octet-stream"
)

func DetermineMimeType(name string) string {
	if ext := path.Ext(name); ext != "" {
		if mt := mime.TypeByExtension(ext); mt != "" {
			return mt
		}
	}
	return DefaultMimeType
}
