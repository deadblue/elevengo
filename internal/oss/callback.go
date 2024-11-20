package oss

import (
	"encoding/json"
	"strings"
)

type Callback struct {
	CallbackUrl  string `json:"callbackUrl"`
	CallbackBody string `json:"callbackBody"`
	// Optional fields
	CallbackHost     *string `json:"callbackHost,omitempty"`
	CallbackSNI      *bool   `json:"callbackSNI,omitempty"`
	CallbackBodyType *string `json:"callbackBodyType,omitempty"`
}

func ReplaceCallbackSha1(callback, fileSha1 string) string {
	cbObj := &Callback{}
	if err := json.Unmarshal([]byte(callback), cbObj); err != nil {
		return callback
	}
	cbObj.CallbackBody = strings.ReplaceAll(cbObj.CallbackBody, "${sha1}", fileSha1)
	cbData, _ := json.Marshal(cbObj)
	return string(cbData)
}
