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
	// Print a warning level log.
	Warn(message string)
	// Print an error level log.
	Error(message string)
}

type simpleLogger struct{}

func (l *simpleLogger) print(level, message string) {
	now := time.Now().Format(time.RFC3339)
	fmt.Printf("%s %s - %s\n", now, level, message)
}

func (l *simpleLogger) Debug(message string) {
	l.print("D", message)
}

func (l *simpleLogger) Info(message string) {
	l.print("I", message)
}

func (l *simpleLogger) Warn(message string) {
	l.print("W", message)
}

func (l *simpleLogger) Error(message string) {
	l.print("E", message)
}

func defaultLogger() Logger {
	return &simpleLogger{}
}
