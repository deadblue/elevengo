package multipart

import (
	"bytes"
	"fmt"
	"io"

	"github.com/deadblue/elevengo/internal/util"
)

type FormBuilder interface {
	AddValue(name, value string) FormBuilder
	AddFile(name, filename string, filesize int64, r io.Reader) FormBuilder
	Build() *Form
}

type implFormBuilder struct {
	// Boundary
	b string
	// Total size
	s int64
	// Readers
	rs []io.Reader
	// Reader number
	rn int
}

func (b *implFormBuilder) buffer() *bytes.Buffer {
	var buf *bytes.Buffer
	if b.rn == 0 {
		// Create bytes.Buffer and append it to rs
		buf = &bytes.Buffer{}
		b.rs = append(b.rs, buf)
		b.rn += 1
	} else {
		// The last reader is always bytes.Buffer!
		buf = (b.rs[b.rn-1]).(*bytes.Buffer)
	}
	return buf
}

func (b *implFormBuilder) incrSize(s int64) {
	if s < 0 {
		b.s = -1
	} else {
		if b.s >= 0 {
			b.s += s
		}
	}
}

func (b *implFormBuilder) AddValue(name string, value string) FormBuilder {
	// Calculate part size
	size := len(b.b) + len(name) + len(value) + 49
	// Write part
	buf := b.buffer()
	buf.Grow(size)
	buf.WriteString("--")
	buf.WriteString(b.b)
	buf.WriteString("\r\nContent-Disposition: form-data; name=\"")
	buf.WriteString(name)
	buf.WriteString("\"\r\n\r\n")
	buf.WriteString(value)
	buf.WriteString("\r\n")
	// Update size
	b.incrSize(int64(size))
	return b
}

func (b *implFormBuilder) AddFile(name string, filename string, filesize int64, r io.Reader) FormBuilder {
	// Calculate part header size
	size := len(b.b) + len(name) + len(filename) + 60
	// Write part header
	buf := b.buffer()
	buf.Grow(size)
	buf.WriteString("--")
	buf.WriteString(b.b)
	buf.WriteString("\r\nContent-Disposition: form-data; name=\"")
	buf.WriteString(name)
	buf.WriteString("\"; filename=\"")
	buf.WriteString(filename)
	buf.WriteString("\"\r\n\r\n")
	b.incrSize(int64(size))

	// Write part tail
	buf = &bytes.Buffer{}
	buf.WriteString("\r\n")

	b.rs = append(b.rs, r, buf)
	b.rn += 2

	// Update form size
	if filesize <= 0 {
		filesize = util.GuessSize(r)
	}
	b.incrSize(filesize + 2)

	return b
}

func (b *implFormBuilder) Build() *Form {
	if b.rn == 0 {
		return nil
	}
	// Write form tail
	size := len(b.b) + 4
	buf := b.buffer()
	buf.Grow(size)
	buf.WriteString("--")
	buf.WriteString(b.b)
	buf.WriteString("--")
	b.s += int64(size)
	// Build form
	return &Form{
		t:  fmt.Sprintf("multipart/form-data; boundary=%s", b.b),
		s:  b.s,
		rs: b.rs,
		ri: 0,
		rn: b.rn,
	}
}

func Builder() FormBuilder {
	return &implFormBuilder{
		b:  generateBoundary("--ElevenGo--", 16),
		s:  0,
		rn: 0,
	}
}
