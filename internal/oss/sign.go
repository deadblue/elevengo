package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/deadblue/elevengo/internal/util"
	"strings"
)

type RequestMetadata struct {
	// Request verb
	Verb string
	// Request header
	Header map[string]string
	// OSS bucket name
	Bucket string
	// OSS object name
	Object string
	// OSS Parameters
	Params map[string]string
}

// CalculateAuthorization calculates authorization for OSS request
func CalculateAuthorization(metadata *RequestMetadata, keyId string, keySecret string) string {
	// Create signer
	signer := hmac.New(sha1.New, []byte(keySecret))
	wx := util.UpgradeWriter(signer)

	// Common parameters
	contentMd5 := metadata.Header[HeaderContentMd5]
	contentType := metadata.Header[HeaderContentType]
	date := metadata.Header[HeaderDate]
	wx.MustWriteString(metadata.Verb, "\n", contentMd5, "\n", contentType, "\n", date, "\n")

	// Canonicalized OSS Headers
	headers := make([]*Pair, 0, len(metadata.Header))
	for name, value := range metadata.Header {
		name = strings.ToLower(name)
		if strings.HasPrefix(name, headerPrefixOss) {
			headers = append(headers, &Pair{
				First: name,
				Last:  value,
			})
		}
	}
	sortPairs(headers)
	for _, header := range headers {
		wx.MustWriteString(header.First, ":", header.Last, "\n")
	}

	// Canonicalized Resource
	wx.MustWriteString("/", metadata.Bucket, "/", metadata.Object)
	// Sub resources
	if metadata.Params != nil && len(metadata.Params) > 0 {
		params := make([]*Pair, 0, len(metadata.Params))
		for name, value := range metadata.Params {
			if _, ok := signingKeyMap[name]; ok {
				params = append(params, &Pair{
					First: name,
					Last:  value,
				})
			}
		}
		sortPairs(params)
		for index, param := range params {
			if index == 0 {
				wx.MustWriteString("?", param.First)
			} else {
				wx.MustWriteString("&", param.First)
			}
			if param.Last != "" {
				wx.MustWriteString("=", param.Last)
			}
		}
	}

	signature := base64.StdEncoding.EncodeToString(signer.Sum(nil))
	return fmt.Sprintf("OSS %s:%s", keyId, signature)
}
