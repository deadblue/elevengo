package plugin

// Logger interface.
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
