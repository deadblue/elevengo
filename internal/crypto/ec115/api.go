package ec115

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"math/rand"

	"github.com/deadblue/elevengo/internal/crypto/lz4"
)

var (
	crcSalt = []byte("^j>WD3Kr?J2gLFjD4W2y@")

	errInvalidEncodedData = errors.New("invalid ec data")
)

func (c *Cipher) EncodeToken(timestamp int64) string {
	buf := make([]byte, 48)
	// Put information to buf
	copy(buf[0:15], c.pubKey[:15])
	copy(buf[24:39], c.pubKey[15:])
	buf[16], buf[40] = 115, 1
	binary.LittleEndian.PutUint32(buf[20:], uint32(timestamp))
	// Encode it
	r1, r2 := byte(rand.Intn(0xff)), byte(rand.Intn(0xff))
	for i := 0; i < 44; i++ {
		if i < 24 {
			buf[i] ^= r1
		} else {
			buf[i] ^= r2
		}
	}
	// Calculate checksum
	h := crc32.NewIEEE()
	h.Write(crcSalt)
	h.Write(buf[:44])
	// Save checksum at the end
	binary.LittleEndian.PutUint32(buf[44:], h.Sum32())
	return base64.StdEncoding.EncodeToString(buf)
}

func (c *Cipher) Encode(input []byte) (output []byte) {
	// Zero padding
	plaintext, plainSize := input, len(input)
	if padSize := aes.BlockSize - (plainSize % aes.BlockSize); padSize != aes.BlockSize {
		plaintext = make([]byte, plainSize+padSize)
		copy(plaintext, input)
		// Make sure all padding bytes are zero
		for i := 0; i < padSize; i++ {
			plaintext[plainSize+i] = 0
		}
		plainSize += padSize
	}
	// Initialize encrypter
	block, _ := aes.NewCipher(c.aesKey)
	enc := cipher.NewCBCEncrypter(block, c.aesIv)
	// Encrypt
	output = make([]byte, plainSize)
	enc.CryptBlocks(output, plaintext)
	return
}

func (c *Cipher) Decode(input []byte) (output []byte, err error) {
	cryptoSize := len(input) - 12
	if cryptoSize < 12 {
		return nil, errInvalidEncodedData
	}
	cryptotext, tail := input[:cryptoSize], input[cryptoSize:]
	// Validate input data
	h := crc32.NewIEEE()
	h.Write(crcSalt)
	h.Write(tail[0:8])
	if h.Sum32() != binary.LittleEndian.Uint32(tail[8:12]) {
		return nil, errInvalidEncodedData
	}
	// Get output size
	for i := 0; i < 4; i++ {
		tail[i] ^= tail[7]
	}
	outputSize := binary.LittleEndian.Uint32(tail[0:4])
	output = make([]byte, outputSize)
	// Initialize decrypter
	block, _ := aes.NewCipher(c.aesKey)
	dec := cipher.NewCBCDecrypter(block, c.aesIv)
	// Decrypt
	plaintext := make([]byte, cryptoSize)
	dec.CryptBlocks(plaintext, cryptotext)
	// Uncompress
	for buf := output; err == nil && outputSize > 0; {
		// Each block is 8192 bytes at maximum
		bufSize := outputSize
		if bufSize > 8192 {
			bufSize = 8192
		}
		srcSize := binary.LittleEndian.Uint16(plaintext[0:2])
		err = lz4.BlockUncompress(plaintext[2:srcSize+2], buf)
		// Prepare for next block
		plaintext = plaintext[srcSize+2:]
		buf = buf[bufSize:]
		outputSize -= bufSize
	}
	return
}
