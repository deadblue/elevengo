package mobile

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/pierrec/lz4/v4"
	"hash/crc32"
	"io"
	"io/ioutil"
	"log"
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
	privKey, x, y, _ := elliptic.GenerateKey(ecCurve, rand.Reader)
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
	le.PutUint32(buf[16:20], c.uid)
	le.PutUint32(buf[20:24], timestamp)
	le.PutUint32(buf[40:44], uint32(apiId))
	// Xor the data
	r1, r2 := randByte(), randByte()
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

func (c *Client) ecDecode(r io.Reader, result interface{}) (err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	dataSize := len(data)
	if dataSize%16 != 12 {
		return errMalformedBody
	}
	body, tail := data[:dataSize-12], data[dataSize-12:]
	bodySize := dataSize - 12
	// Get information from tail
	dataSize, encrypted, compressed, err := parseTail(tail)
	if err != nil {
		return
	}
	// Decrypt
	if encrypted {
		block, _ := aes.NewCipher(c.aesKey)
		dec := cipher.NewCBCDecrypter(block, c.aesIv)
		plain := make([]byte, bodySize)
		dec.CryptBlocks(plain, body)
		// De-padding
		for j := bodySize - 1; ; j-- {
			if plain[j] == 0 {
				plain = plain[:j]
			} else {
				break
			}
		}
		body = plain
	}
	// Decompress
	if compressed {
		buf := make([]byte, dataSize)
		compSize := le.Uint16(body[:2])
		if _, err = lz4.UncompressBlock(body[2:compSize+2], buf); err == nil {
			body = buf
		} else {
			return
		}
	}
	log.Printf("Body: %s", body)
	return json.Unmarshal(body, result)
}

func parseTail(tail []byte) (size int, encrypted, compressed bool, err error) {
	// Check CRC32
	crc := crc32.NewIEEE()
	_, _ = crc.Write(crcSalt)
	_, _ = crc.Write(tail[:8])
	if crc.Sum32() != le.Uint32(tail[8:]) {
		err = errMalformedBody
		return
	}
	// Flags
	compressed = tail[4] == 0x01
	encrypted = tail[5] == 0x01
	// Original data size
	key := tail[7]
	for i := 0; i < 4; i++ {
		tail[i] = tail[i] ^ key
	}
	size = int(le.Uint32(tail[0:4]))
	return
}

func (c *Client) ecEncode(data []byte) (r io.Reader) {
	// Padding
	plain, size := data, len(data)
	if m := size % 16; m != 0 {
		n := 16 - m
		padding := bytes.Repeat([]byte{0}, n)
		plain = append(plain, padding...)
		size += n
	}
	// Encrypt
	block, _ := aes.NewCipher(c.aesKey)
	enc := cipher.NewCBCEncrypter(block, c.aesIv)
	buf := make([]byte, size)
	enc.CryptBlocks(buf, plain)
	//
	return bytes.NewReader(buf)
}
