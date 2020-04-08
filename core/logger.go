package core

import (
	"fmt"
	"github.com/deadblue/elevengo/plugin"
)

type LoggerEx interface {
	plugin.Logger
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type implLoggerEx struct {
	logger plugin.Logger
}

func (l *implLoggerEx) Debug(message string) {
	if l.logger != nil {
		l.logger.Debug(message)
	}
}
func (l *implLoggerEx) Debugf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Debug(fmt.Sprintf(format, args...))
	}
}
func (l *implLoggerEx) Info(message string) {
	if l.logger != nil {
		l.logger.Info(message)
	}
}
func (l *implLoggerEx) Infof(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Info(fmt.Sprintf(format, args...))
	}
}
func (l *implLoggerEx) Warn(message string) {
	if l.logger != nil {
		l.logger.Warn(message)
	}
}
func (l *implLoggerEx) Warnf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Warn(fmt.Sprintf(format, args...))
	}
}
func (l *implLoggerEx) Error(message string) {
	if l.logger != nil {
		l.logger.Error(message)
	}
}
func (l *implLoggerEx) Errorf(format string, args ...interface{}) {
	if l.logger != nil {
		l.logger.Error(fmt.Sprintf(format, args...))
	}
}

// Wrap plugin.Logger by adding format-able logging methods.
func WrapLogger(logger plugin.Logger) LoggerEx {
	return &implLoggerEx{logger: logger}
}
