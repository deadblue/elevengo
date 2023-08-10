package util

import "strings"

func MarshalCookies(cookies map[string]string) string {
	if len(cookies) > 0 {
		buf, isFirst := strings.Builder{}, true
		for ck, cv := range cookies {
			if !isFirst {
				buf.WriteString("; ")
			}
			buf.WriteString(ck)
			buf.WriteRune('=')
			buf.WriteString(cv)
			isFirst = false
		}
		return buf.String()
	}
	return ""
}
