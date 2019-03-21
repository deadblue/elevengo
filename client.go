package elevengo

import (
	"net/http"
)

type _UserInfo struct {
	UserId string
}

type _OfflineToken struct {
	Sign string
	Time int64
}

type Client struct {
	// basic properties for communication
	userAgent string

	jar    http.CookieJar
	client *http.Client

	info    *_UserInfo
	offline *_OfflineToken
}

func New() (client *Client, err error) {
	return Custom(defaultUserAgent)
}

func Custom(ua string) (client *Client, err error) {
	client = &Client{userAgent: ua}
	if err = client.setup(); err != nil {
		client = nil
	}
	return
}
