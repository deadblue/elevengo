package sso

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func hexEncode(n int64, length int) string {
	s := strconv.FormatInt(n, 16)
	if l := len(s); l < length {
		s = strings.Repeat("0", length-l) + s
	} else if l > length {
		s = s[l-length:]
	}
	return s
}

func GenerateExt() string {
	return hexEncode(time.Now().Unix(), 8) +
		hexEncode(rand.Int63n(123456789), 5)
}
