package ec115

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"filippo.io/nistec"
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

func New() *Cipher {
	curve := elliptic.P224()
	// Generate local key
	localKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	scalar := make([]byte, (curve.Params().BitSize+7)/8)
	localKey.D.FillBytes(scalar)
	pubKey := elliptic.MarshalCompressed(curve, localKey.X, localKey.Y)
	// Parse remote key
	remoteKey, _ := nistec.NewP224Point().SetBytes(serverKey)
	// ECDH key exchange
	sharedPoint, _ := nistec.NewP224Point().ScalarMult(remoteKey, scalar)
	sharedSecret, _ := sharedPoint.BytesX()
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
