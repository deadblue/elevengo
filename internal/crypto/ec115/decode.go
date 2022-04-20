package ec115

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/pierrec/lz4/v4"
	"hash/crc32"
	"io"
	"io/ioutil"
)

func (c *Coder) DecodeData(r io.Reader) (result []byte, err error) {
	// Read all data from reader
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	// Check data size
	dataSize := len(data)
	if dataSize%16 != 12 {
		err = errMalformedData
		return
	}
	// Get information from tail
	realSize, encrypted, compressed, err := parseTail(data[dataSize-12:])
	if err != nil {
		return
	}
	result, bodySize := data[:dataSize-12], dataSize-12
	// Decrypt
	if encrypted {
		block, _ := aes.NewCipher(c.aesKey)
		dec := cipher.NewCBCDecrypter(block, c.aesIv)
		plain := make([]byte, bodySize)
		dec.CryptBlocks(plain, result)
		// De-padding
		for j := bodySize - 1; plain[j] == 0; j-- {
			plain = plain[:j]
		}
		result = plain
	}
	// Decompress
	if compressed {
		buf := make([]byte, realSize)
		compSize := le.Uint16(result[:2])
		if _, err = lz4.UncompressBlock(result[2:compSize+2], buf); err == nil {
			result = buf
		}
	}
	return
}

func parseTail(tail []byte) (size int, encrypted, compressed bool, err error) {
	// Check CRC32
	crc := crc32.NewIEEE()
	_, _ = crc.Write(salt)
	_, _ = crc.Write(tail[:8])
	if crc.Sum32() != le.Uint32(tail[8:]) {
		err = errMalformedData
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
