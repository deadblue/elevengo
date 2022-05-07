package option

import "github.com/deadblue/elevengo/plugin"

type HttpOption struct {
	Client plugin.HttpClient
}

func (o *HttpOption) isOption() {}
