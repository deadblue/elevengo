package util

import "strings"

func SecretUrl(url string) string {
	if index := strings.IndexRune(url, ':'); index > 0 {
		scheme := strings.ToLower(url[:index])
		if scheme == "http" {
			return "https" + url[index:]
		}
	}
	return url
}
