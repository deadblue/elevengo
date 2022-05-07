package ec115

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"hash/crc32"
	"io"
	"time"
)

func (c *Coder) EncodeKey(userId uint32) string {
	timestamp := uint32(time.Now().Unix())
	// Fill buffer
	buf := make([]byte, 48)
	buf[15], buf[39] = 0, 0
	copy(buf[0:15], c.pubKey[0:15])
	copy(buf[24:39], c.pubKey[15:30])
	le.PutUint32(buf[16:20], userId)
	le.PutUint32(buf[20:24], timestamp)
	// Update counter
	le.PutUint32(buf[40:44], c.counter)
	c.counter += 1
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
	_, _ = crc.Write(salt)
	_, _ = crc.Write(buf[:44])
	le.PutUint32(buf[44:48], crc.Sum32())
	// Encoding to base64
	return base64.StdEncoding.EncodeToString(buf)
}

func (c *Coder) EncodeData(data []byte) io.Reader {
	// Zero-Padding
	plain, size := data, len(data)
	if m := size % 16; m != 0 {
		n := 16 - m
		padding := bytes.Repeat([]byte{0}, n)
		plain = append(plain, padding...)
		size += n
	}
	// AES Encrypt
	block, _ := aes.NewCipher(c.aesKey)
	enc := cipher.NewCBCEncrypter(block, c.aesIv)
	buf := make([]byte, size)
	enc.CryptBlocks(buf, plain)
	// Wrap to reader
	return bytes.NewReader(buf)
}
