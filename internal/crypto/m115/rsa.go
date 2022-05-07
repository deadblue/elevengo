package m115

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

var (
	rsaPrivateKey = []byte("-----BEGIN RSA PRIVATE KEY-----\n" +
		"MIICXAIBAAKBgQCMgUJLwWb0kYdW6feyLvqgNHmwgeYYlocst8UckQ1+waTOKHFC\n" +
		"TVyRSb1eCKJZWaGa08mB5lEu/asruNo/HjFcKUvRF6n7nYzo5jO0li4IfGKdxso6\n" +
		"FJIUtAke8rA2PLOubH7nAjd/BV7TzZP2w0IlanZVS76n8gNDe75l8tonQQIDAQAB\n" +
		"AoGANwTasA2Awl5GT/t4WhbZX2iNClgjgRdYwWMI1aHbVfqADZZ6m0rt55qng63/\n" +
		"3NsjVByAuNQ2kB8XKxzMoZCyJNvnd78YuW3Zowqs6HgDUHk6T5CmRad0fvaVYi6t\n" +
		"viOkxtiPIuh4QrQ7NUhsLRtbH6d9s1KLCRDKhO23pGr9vtECQQDpjKYssF+kq9iy\n" +
		"A9WvXRjbY9+ca27YfarD9WVzWS2rFg8MsCbvCo9ebXcmju44QhCghQFIVXuebQ7Q\n" +
		"pydvqF0lAkEAmgLnib1XonYOxjVJM2jqy5zEGe6vzg8aSwKCYec14iiJKmEYcP4z\n" +
		"DSRms43hnQsp8M2ynjnsYCjyiegg+AZ87QJANuwwmAnSNDOFfjeQpPDLy6wtBeft\n" +
		"5VOIORUYiovKRZWmbGFwhn6BQL+VaafrNaezqUweBRi1PYiAF2l3yLZbUQJAf/nN\n" +
		"4Hz/pzYmzLlWnGugP5WCtnHKkJWoKZBqO2RfOBCq+hY4sxvn3BHVbXqGcXLnZPvo\n" +
		"YuaK7tTXxZSoYLEzeQJBAL8Mt3AkF1Gci5HOug6jT4s4Z+qDDrUXo9BlTwSWP90v\n" +
		"wlHF+mkTJpKd5Wacef0vV+xumqNorvLpIXWKwxNaoHM=\n" +
		"-----END RSA PRIVATE KEY-----")
	rsaPublicKey = []byte("-----BEGIN RSA PUBLIC KEY-----\n" +
		"MIGJAoGBANHetaZ5idEKXAsEHRGrR2Wbwys+ZakvkjbdLMIUCg2klfoOfvh19vrL\n" +
		"TZgfXl47peZ4Ed1zt6QQUlQiL6zCBqdOiREhVFGv/PXr/eiHvJrbZ1wCqDX3XL53\n" +
		"pgOvggaD9DnnztQokyPfnJBVdp4VeYuUU+iQWLPi4/GGsHsEapltAgMBAAE=\n" +
		"-----END RSA PUBLIC KEY-----")

	/* Client Key */
	rsaClientKey *rsa.PrivateKey
	/* Server Key */
	rsaServerKey *rsa.PublicKey
)

func rsaEncrypt(input []byte) []byte {
	plainSize, blockSize := len(input), rsaServerKey.Size()-11
	buf := bytes.Buffer{}
	for offset := 0; offset < plainSize; offset += blockSize {
		sliceSize := blockSize
		if offset+sliceSize > plainSize {
			sliceSize = plainSize - offset
		}
		slice, _ := rsa.EncryptPKCS1v15(
			rand.Reader, rsaServerKey, input[offset:offset+sliceSize])
		buf.Write(slice)
	}
	return buf.Bytes()
}

func rsaDecrypt(input []byte) []byte {
	output := make([]byte, 0)
	cipherSize, blockSize := len(input), rsaServerKey.Size()
	for offset := 0; offset < cipherSize; offset += blockSize {
		sliceSize := blockSize
		if offset+sliceSize > cipherSize {
			sliceSize = cipherSize - offset
		}
		slice, _ := rsa.DecryptPKCS1v15(
			rand.Reader, rsaClientKey, input[offset:offset+sliceSize])
		output = append(output, slice...)
	}
	return output
}

func init() {
	// Parse client private key
	block, _ := pem.Decode(rsaPrivateKey)
	rsaClientKey, _ = x509.ParsePKCS1PrivateKey(block.Bytes)
	// Parse server public key
	block, _ = pem.Decode(rsaPublicKey)
	rsaServerKey, _ = x509.ParsePKCS1PublicKey(block.Bytes)
}
