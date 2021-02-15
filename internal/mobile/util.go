package mobile

import "math/rand"

func randByte() byte {
	return byte(rand.Intn(0xff))
}
