package ec115

import (
	"crypto/elliptic"
	"crypto/rand"
)

var (
	serverKey = []byte{
		0x04, 0x57, 0xa2, 0x92, 0x57, 0xcd, 0x23, 0x20,
		0xe5, 0xd6, 0xd1, 0x43, 0x32, 0x2f, 0xa4, 0xbb,
		0x8a, 0x3c, 0xf9, 0xd3, 0xcc, 0x62, 0x3e, 0xf5,
		0xed, 0xac, 0x62, 0xb7, 0x67, 0x8a, 0x89, 0xc9,
		0x1a, 0x83, 0xba, 0x80, 0x0d, 0x61, 0x29, 0xf5,
		0x22, 0xd0, 0x34, 0xc8, 0x95, 0xdd, 0x24, 0x65,
		0x24, 0x3a, 0xdd, 0xc2, 0x50, 0x95, 0x3b, 0xee,
		0xba,
	}
)

type Cipher struct {
	// Client public key
	pubKey []byte
	// AES key & IV
	aesKey []byte
	aesIv  []byte
}

// WARNING:
// Some elliptic APIs are deprecated since 1.21.0, but the replacing APIs do not
// support P224.
func New() *Cipher {
	curve := elliptic.P224()
	// Generate client key-pair
	privKey, x, y, _ := elliptic.GenerateKey(
		curve, rand.Reader,
	)
	pubKey := elliptic.MarshalCompressed(curve, x, y)
	// Parse server key
	serverX, serverY := elliptic.Unmarshal(curve, serverKey)
	// ECDH key exchange
	sharedX, _ := curve.ScalarMult(serverX, serverY, privKey)
	sharedSecret := sharedX.Bytes()
	// Instantiate cipher
	pubKeySize, sharedSecretSize := len(pubKey), len(sharedSecret)
	cipher := &Cipher{
		pubKey: make([]byte, pubKeySize+1),
		aesKey: make([]byte, 16),
		aesIv:  make([]byte, 16),
	}
	cipher.pubKey[0] = byte(pubKeySize)
	copy(cipher.pubKey[1:], pubKey)
	copy(cipher.aesKey, sharedSecret[:16])
	copy(cipher.aesIv, sharedSecret[sharedSecretSize-16:])
	return cipher
}
