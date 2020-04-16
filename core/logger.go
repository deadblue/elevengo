package core

import (
	"fmt"
	"github.com/deadblue/elevengo/plugin"
)

type LoggerEx interface {
	plugin.Logger
	Printf(format string, v ...interface{})
}

type implLoggerEx struct {
	logger plugin.Logger
}

func (i *implLoggerEx) Println(v ...interface{}) {
	if i.logger == nil {
		return
	}
	i.logger.Println(v...)
}
func (i *implLoggerEx) Printf(format string, v ...interface{}) {
	if i.logger == nil {
		return
	}
	i.logger.Println(fmt.Sprintf(format, v...))
}

// Extend the logger.
func wrapEx(logger plugin.Logger) LoggerEx {
	if lx, ok := logger.(LoggerEx); ok {
		return lx
	} else {
		return &implLoggerEx{logger: logger}
	}
}
