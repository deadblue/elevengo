package mobile

import (
	"crypto/rsa"
	"net/http"
)

const (
	AppVersion = "26.1.0"
)

type Client struct {
	hc *http.Client
	cj http.CookieJar

	// User ID
	userId uint32

	// EC public key
	ecPubKey []byte
	// AES key
	aesKey []byte
	aesIv  []byte
	// RSA keys
	rsaPrivKey *rsa.PrivateKey
	rsaPubKey  *rsa.PublicKey
}

func New() (client *Client, err error) {
	client = &Client{}
	client.initHttpClient()
	client.ecInit()
	client.rasInit()

	return
}
