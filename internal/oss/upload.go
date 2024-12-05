package oss

import (
	"fmt"
)

const (
	Region = "cn-shenzhen"

	endpointHost = "oss-cn-shenzhen.aliyuncs.com"
)

func GetPutObjectUrl(bucket, key string) string {
	return fmt.Sprintf("https://%s.%s/%s", bucket, endpointHost, key)
}

func GetEndpointUrl() string {
	return fmt.Sprintf("https://%s", endpointHost)
}
