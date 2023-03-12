package hash

import (
	"encoding/base64"
	"encoding/hex"
	"hash"
	"strings"
)

func ToHex(h hash.Hash) string {
	return hex.EncodeToString(h.Sum(nil))
}

func ToHexUpper(h hash.Hash) string {
	return strings.ToUpper(ToHex(h))
}

func ToBase64(h hash.Hash) string {
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
