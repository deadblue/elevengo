package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Hex(str string) string {
	h := md5.Sum([]byte(str))
	return hex.EncodeToString(h[:])
}
