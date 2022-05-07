package multipart

import (
	"io"
)

type Form struct {
	// Content type
	t string
	// Total Size
	s int64
	// Readers
	rs []io.Reader
	// Reader index & count
	ri, rn int
}

func (f *Form) ContentType() string {
	return f.t
}

func (f *Form) ContentLength() int64 {
	return f.s
}

func (f *Form) Read(p []byte) (n int, err error) {
	err = io.EOF
	for f.ri < f.rn && err == io.EOF {
		r, n0 := f.rs[f.ri], 0

		n0, err = r.Read(p)
		n += n0
		if err == io.EOF {
			f.ri += 1
			p = p[n0:]
		}
	}
	return
}
