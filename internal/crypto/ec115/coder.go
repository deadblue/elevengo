package ec115

type Coder struct {
	pubKey []byte
	aesKey []byte
	aesIv  []byte

	counter uint32
}

func New() *Coder {
	return (&Coder{
		pubKey:  make([]byte, 30, 30),
		aesKey:  make([]byte, 16, 16),
		aesIv:   make([]byte, 16, 16),
		counter: 0,
	}).init()
}
