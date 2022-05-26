package hash

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
)

func Sha1HexUpper(r io.Reader) (digest string) {
	h := sha1.New()
	_, _ = io.Copy(h, r)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
