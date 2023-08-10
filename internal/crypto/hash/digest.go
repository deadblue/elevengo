package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
)

type DigestResult struct {
	Size int64
	SHA1 string
	MD5  string
}

func Digest(r io.Reader, result *DigestResult) (err error) {
	hs, hm := sha1.New(), md5.New()
	w := io.MultiWriter(hs, hm)
	// Write remain data.
	if result.Size, err = io.Copy(w, r); err != nil {
		return
	}
	result.SHA1, result.MD5 = ToHexUpper(hs), ToBase64(hm)
	return nil
}

func DigestRange(r io.ReadSeeker, rangeSpec string) (result string, err error) {
	var start, end int64
	if _, err = fmt.Sscanf(rangeSpec, "%d-%d", &start, &end); err != nil {
		return
	}
	h := sha1.New()
	r.Seek(start, io.SeekStart)
	if _, err = io.CopyN(h, r, end-start+1); err == nil {
		result = ToHexUpper(h)
	}
	return
}
