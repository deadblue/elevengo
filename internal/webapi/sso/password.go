package sso

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func sha1Hex(input string) string {
	h := sha1.Sum([]byte(input))
	return hex.EncodeToString(h[:])
}

func EncodePassword(account, password, ext string) string {
	return sha1Hex(sha1Hex(sha1Hex(password)+sha1Hex(account)) + strings.ToUpper(ext))
}

func getRuneWeight(r rune) int {
	if r >= '0' && r <= '9' {
		return 1
	} else if r >= 'A' && r <= 'Z' {
		return 2
	} else if r >= 'a' && r <= 'z' {
		return 4
	} else {
		return 8
	}
}

func GetPasswordLevel(password string) int {
	n := len(password)
	if n < 5 {
		return 0
	}
	weight := 0
	for _, r := range password {
		weight |= getRuneWeight(r)
	}
	level := 0
	for ; weight != 0; weight >>= 1 {
		if weight&1 == 1 {
			level += 1
		}
	}
	if n > 8 {
		level += 1
	}
	return level
}
