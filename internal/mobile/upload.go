package mobile

import (
	"log"
)

type _UploadUserKey struct {
	UserKey string `json:"userkey"`
}

func (c *Client) getUploadKey() (err error) {
	result := make(map[string]interface{})
	err = c.callApi(apiUploadGetKey, map[string]string{
		"app_id": "0",
	}, nil, &result)
	log.Printf("Result => %#v", result)
	return
}
