package lz4

import (
	"errors"
)

var (
	errOutOfRange = errors.New("out of range")
)

type Buffer struct {
	buf    []byte
	cursor int
	size   int
}

func newBuffer(p []byte) *Buffer {
	return &Buffer{
		buf:    p,
		cursor: 0,
		size:   len(p),
	}
}

func (b *Buffer) AtTheEnd() bool {
	return b.cursor == b.size
}

// Next return a slice of underlying buf for reading/writing.
func (b *Buffer) Next(n int) (p []byte, err error) {
	if b.cursor+n > b.size {
		err = errOutOfRange
	} else {
		p = b.buf[b.cursor : b.cursor+n]
		b.cursor += n
	}
	return
}

func (b *Buffer) ReadByte() (v byte, err error) {
	if b.cursor == b.size {
		err = errOutOfRange
	} else {
		v = b.buf[b.cursor]
		b.cursor += 1
	}
	return
}

func (b *Buffer) ReadToken() (literalLen, matchLen int, err error) {
	var v byte
	if v, err = b.ReadByte(); err == nil {
		literalLen = int(v >> 4)
		matchLen = 4 + int(v&0x0f)
	}
	return
}

func (b *Buffer) ReadExtraLength() (l int, err error) {
	var v byte
	for {
		if v, err = b.ReadByte(); err == nil {
			l += int(v)
		}
		if err != nil || v != 0xff {
			break
		}
	}
	return
}

func (b *Buffer) ReadOffset() (offset int, err error) {
	var p []byte
	if p, err = b.Next(2); err == nil {
		offset = int(p[0]) | (int(p[1]) << 8)
	}
	return
}

func (b *Buffer) WriteMatchedLiteral(offset, length int) (err error) {
	start := b.cursor - offset
	var p []byte
	if p, err = b.Next(length); err != nil {
		return
	}
	copy(p, b.buf[start:])
	return
}

func CopyN(dst, src *Buffer, n int) (err error) {
	var sp, dp []byte
	if sp, err = src.Next(n); err != nil {
		return
	}
	if dp, err = dst.Next(n); err != nil {
		return
	}
	copy(dp, sp)
	return
}
