package util

import (
	"bytes"
	"io"
	"os"
	"strings"
)

type WriterEx struct {
	w io.Writer
}

func (w *WriterEx) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *WriterEx) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

func (w *WriterEx) WriteByte(b byte) (err error) {
	_, err = w.Write([]byte{b})
	return
}

func (w *WriterEx) MustWriteString(s ...string) {
	for _, item := range s {
		if _, err := w.WriteString(item); err != nil {
			break
		}
	}
}

// UpgradeWriter gives you a powerful Writer than the original one!
func UpgradeWriter(w io.Writer) *WriterEx {
	return &WriterEx{w: w}
}

func ConsumeReader(r io.ReadCloser) {
	_, _ = io.Copy(io.Discard, r)
	_ = r.Close()
}

func QuietlyClose(c io.Closer) {
	_ = c.Close()
}

func GuessSize(r io.Reader) (size int64) {
	size = -1
	switch r := r.(type) {
	case *bytes.Buffer:
		size = int64(r.Len())
	case *bytes.Reader:
		size = r.Size()
	case *strings.Reader:
		size = int64(r.Len())
	case *os.File:
		if i, e := r.Stat(); e == nil {
			size = i.Size()
		}
	}
	return
}
