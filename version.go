package elevengo

import (
	"fmt"
	"runtime"
)

const (
	libName = "elevengo"
	libVer  = "0.6.1"
)

var (
	version = fmt.Sprintf("%s %s (%s %s/%s)",
		libName, libVer, runtime.Version(), runtime.GOOS, runtime.GOARCH)
)

func (a *Agent) Version() string {
	return version
}
