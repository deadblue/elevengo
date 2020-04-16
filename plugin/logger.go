package plugin

/*
Logger interface for printing debug message.

Caller can implement himself, or simply uses log.Logger in stdlib.
*/
type Logger interface {

	// Print debug message.
	// NOTICE: The message does not end with newline character("\n"),
	// implementation should append one.
	Println(v ...interface{})
}
