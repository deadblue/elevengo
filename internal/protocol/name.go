package protocol

import "fmt"

const (
	namePrefix = "Mozilla/5.0"
)

func MakeUserAgent(name, appName, appVer string) string {
	if name == "" {
		return fmt.Sprintf("%s %s/%s", namePrefix, appName, appVer)
	} else {
		return fmt.Sprintf("%s %s %s/%s", namePrefix, name, appName, appVer)
	}
}
