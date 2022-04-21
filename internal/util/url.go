package util

import "strings"

func SecretUrl(url string) string {
	index := strings.IndexRune(url, ':')
	if index < 0 {
		return url
	}
	scheme := strings.ToLower(url[:index])
	if scheme == "http" {
		return "https" + url[index:]
	}
	return url
}
