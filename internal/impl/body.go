package impl

import "io"

type _BodyImpl struct {
	rc io.ReadCloser

	size int64
}

func (i *_BodyImpl) Read(p []byte) (int, error) {
	return i.rc.Read(p)
}

func (i *_BodyImpl) Close() error {
	return i.rc.Close()
}

func (i *_BodyImpl) Size() int64 {
	return i.size
}
