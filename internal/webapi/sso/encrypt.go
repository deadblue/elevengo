package sso

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

var (
	errInvalidKey = errors.New("invalid key")
)

func parsePublicKey(pemData []byte) (pubKey *rsa.PublicKey, err error) {
	b, _ := pem.Decode(pemData)
	if b == nil {
		return nil, errInvalidKey
	}
	key, ok := any(nil), false
	switch b.Type {
	case "PUBLIC KEY":
		key, err = x509.ParsePKIXPublicKey(b.Bytes)
		if err == nil {
			if pubKey, ok = key.(*rsa.PublicKey); !ok {
				err = errInvalidKey
			}
		}
	case "RSA PUBLIC KEY":
		pubKey, err = x509.ParsePKCS1PublicKey(b.Bytes)
	default:
		err = errInvalidKey
	}
	return
}

func EncryptPassword(password string, time int64, key string) (pwd string, err error) {
	pubKey, err := parsePublicKey([]byte(key))
	if err != nil {
		return
	}
	input := fmt.Sprintf("%s_%d", sha1Hex(password), time)
	output, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, []byte(input))
	if err == nil {
		pwd = base64.StdEncoding.EncodeToString(output)
	}
	return
}
