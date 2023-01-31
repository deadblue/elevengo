package oss

import (
	"fmt"
)

const (
	endpoint = "oss-cn-shenzhen.aliyuncs.com"
)

func GetPutObjectUrl(bucket, key string) string {
	return fmt.Sprintf("https://%s.%s/%s", bucket, endpoint, key)
}

func GetEndpointUrl() string {
	return fmt.Sprintf("https://%s", endpoint)
}
