package oss

import (
	"bytes"
	"sort"
)

const (
	headerPrefixOss = "x-oss-"
)

var (
	signingKeyMap = map[string]bool{
		"acl":                          true,
		"uploads":                      true,
		"location":                     true,
		"cors":                         true,
		"logging":                      true,
		"website":                      true,
		"referer":                      true,
		"lifecycle":                    true,
		"delete":                       true,
		"append":                       true,
		"tagging":                      true,
		"objectMeta":                   true,
		"uploadId":                     true,
		"partNumber":                   true,
		"security-token":               true,
		"position":                     true,
		"img":                          true,
		"style":                        true,
		"styleName":                    true,
		"replication":                  true,
		"replicationProgress":          true,
		"replicationLocation":          true,
		"cname":                        true,
		"bucketInfo":                   true,
		"comp":                         true,
		"qos":                          true,
		"live":                         true,
		"status":                       true,
		"vod":                          true,
		"startTime":                    true,
		"endTime":                      true,
		"symlink":                      true,
		"x-oss-process":                true,
		"response-content-type":        true,
		"x-oss-traffic-limit":          true,
		"response-content-language":    true,
		"response-expires":             true,
		"response-cache-control":       true,
		"response-content-disposition": true,
		"response-content-encoding":    true,
		"udf":                          true,
		"udfName":                      true,
		"udfImage":                     true,
		"udfId":                        true,
		"udfImageDesc":                 true,
		"udfApplication":               true,
		"udfApplicationLog":            true,
		"restore":                      true,
		"callback":                     true,
		"callback-var":                 true,
		"qosInfo":                      true,
		"policy":                       true,
		"stat":                         true,
		"encryption":                   true,
		"versions":                     true,
		"versioning":                   true,
		"versionId":                    true,
		"requestPayment":               true,
		"x-oss-request-payer":          true,
		"sequential":                   true,
		"inventory":                    true,
		"inventoryId":                  true,
		"continuation-token":           true,
		"asyncFetch":                   true,
		"worm":                         true,
		"wormId":                       true,
		"wormExtend":                   true,
		"withHashContext":              true,
		"x-oss-enable-md5":             true,
		"x-oss-enable-sha1":            true,
		"x-oss-enable-sha256":          true,
		"x-oss-hash-ctx":               true,
		"x-oss-md5-ctx":                true,
		"transferAcceleration":         true,
		"regionList":                   true,
	}
)

type Pair struct {
	First string
	Last  string
}

type Pairs []*Pair

func (h Pairs) Len() int {
	return len(h)
}

func (h Pairs) Less(i, j int) bool {
	return bytes.Compare([]byte(h[i].First), []byte(h[j].First)) < 0
}

func (h Pairs) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func sortPairs(pairs []*Pair) {
	sort.Sort(Pairs(pairs))
}
