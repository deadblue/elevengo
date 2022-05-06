package multipart

import (
	"crypto/rand"
	"strings"
)

var (
	chars      = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	charsCount = len(chars)
)

func generateBoundary(prefix string, length int) string {
	seed := make([]byte, length)
	_, _ = rand.Read(seed)

	buf, n := strings.Builder{}, len(prefix)
	buf.Grow(length + n)
	if n > 0 {
		buf.WriteString(prefix)
	}
	for _, b := range seed {
		index := int(b) % charsCount
		buf.WriteRune(chars[index])
	}
	return buf.String()
}
