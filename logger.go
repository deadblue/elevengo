package elevengo

import (
	"fmt"
	"time"
)

// Logger interface
type Logger interface {
	// Print a debug level log.
	Debug(message string)
	// Print an information level log.
	Info(message string)
	// Print a warnning level log.
	Warn(message string)
	// Print an error level log.
	Error(message string)
}

// TODO: Inject Logger into Agent.

type defaultLogger struct{}

func (l *defaultLogger) print(level, message string) {
	now := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] %s - %s\n", now, level, message)
}

func (l *defaultLogger) Debug(message string) {
	l.print("D", message)
}

func (l *defaultLogger) Info(message string) {
	l.print("I", message)
}

func (l *defaultLogger) Warn(message string) {
	l.print("W", message)
}

func (l *defaultLogger) Error(message string) {
	l.print("E", message)
}
