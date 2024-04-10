package elevengo

import (
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/types"
)

// LowlevelClient returns low-level client that can be used to directly call ApiSpec.
func (a *Agent) LowlevelClient() client.Client {
	return a.llc
}

// LowlevelParams returns common parameters for low-level API calling.
func (a *Agent) LowlevelParams() types.CommonParams {
	return a.common
}
