package mobile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"github.com/pierrec/lz4/v4"
	"hash/crc32"
	"log"
	"math/rand"
	"time"
)

var (
	ecCurve = elliptic.P224()

	ecSvrX, ecSvrY = elliptic.Unmarshal(ecCurve, []byte{
		0x04, 0x57, 0xa2, 0x92, 0x57, 0xcd, 0x23, 0x20,
		0xe5, 0xd6, 0xd1, 0x43, 0x32, 0x2f, 0xa4, 0xbb,
		0x8a, 0x3c, 0xf9, 0xd3, 0xcc, 0x62, 0x3e, 0xf5,
		0xed, 0xac, 0x62, 0xb7, 0x67, 0x8a, 0x89, 0xc9,
		0x1a, 0x83, 0xba, 0x80, 0x0d, 0x61, 0x29, 0xf5,
		0x22, 0xd0, 0x34, 0xc8, 0x95, 0xdd, 0x24, 0x65,
		0x24, 0x3a, 0xdd, 0xc2, 0x50, 0x95, 0x3b, 0xee,
		0xba,
	})

	crcSalt = []byte("^j>WD3Kr?J2gLFjD4W2y@")

	le = binary.LittleEndian

	errMalformedBody = errors.New("malformed response body")
)

func (c *Client) ecInit() {
	// Generate EC P224 key pair
	privKey, x, y, _ := elliptic.GenerateKey(ecCurve, crand.Reader)
	pubKey := elliptic.MarshalCompressed(ecCurve, x, y)
	// Store public key
	keySize := len(pubKey)
	c.ecPubKey = make([]byte, keySize+1)
	c.ecPubKey[0] = byte(keySize)
	copy(c.ecPubKey[1:], pubKey)
	// ECDH key exchanging
	x, _ = ecCurve.ScalarMult(ecSvrX, ecSvrY, privKey)
	secret := x.Bytes()
	c.aesKey = secret[0:16]
	c.aesIv = secret[len(secret)-16:]
	return
}

func (c *Client) ecEncodeKey(apiId int) string {
	buf := make([]byte, 48)
	timestamp := uint32(time.Now().Unix())
	// Fill buffer
	buf[15], buf[39] = 0, 0
	copy(buf[0:15], c.ecPubKey[0:15])
	copy(buf[24:39], c.ecPubKey[15:30])
	le.PutUint32(buf[16:20], c.userId)
	le.PutUint32(buf[20:24], timestamp)
	le.PutUint32(buf[40:44], uint32(apiId))
	// Xor the data
	r1, r2 := byte(rand.Intn(0xff)), byte(rand.Intn(0xff))
	for i := 0; i < 44; i++ {
		if i < 24 {
			buf[i] = buf[i] ^ r1
		} else {
			buf[i] = buf[i] ^ r2
		}
	}
	// Calculate checksum
	crc := crc32.NewIEEE()
	_, _ = crc.Write(crcSalt)
	_, _ = crc.Write(buf[:44])
	le.PutUint32(buf[44:48], crc.Sum32())
	// Encoding to base64
	return base64.StdEncoding.EncodeToString(buf)
}

func (c *Client) ecDecode(data []byte) (err error) {
	dataSize := len(data)
	if dataSize%16 != 12 {
		return errMalformedBody
	}
	body, tail := data[:dataSize-12], data[dataSize-12:]
	// Verify checksum
	crc := crc32.NewIEEE()
	_, _ = crc.Write(crcSalt)
	_, _ = crc.Write(tail[:8])
	if crc.Sum32() != le.Uint32(tail[8:]) {
		return errMalformedBody
	}
	// Decrypt
	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return
	}
	dec := cipher.NewCBCDecrypter(block, c.aesIv)
	plain := make([]byte, dataSize-12)
	dec.CryptBlocks(plain, body)
	for j := dataSize - 13; ; j-- {
		if plain[j] == 0 {
			plain = plain[:j]
		} else {
			break
		}
	}
	// Decompress
	dataSize = int(le.Uint16(plain[:2]))
	log.Printf("Compress data size: %d", dataSize)

	body = make([]byte, dataSize*2)
	dataSize, err = lz4.UncompressBlock(plain[2:dataSize+2], body)
	if err != nil {
		return err
	}
	log.Printf("Decompress data size: %d", dataSize)
	log.Printf("Body: %s", body[:dataSize])
	return
}
