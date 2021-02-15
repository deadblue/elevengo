package mobile

import (
	"crypto/rsa"
	"net/http"
)

type Client struct {
	// HTTP client
	hc *http.Client
	// Cookie jar
	cj http.CookieJar
	// User ID
	uid uint32
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
	client.httpInit()
	client.ecInit()

	// TODO: Move rsa keys out of client
	client.rasInit()

	return
}
