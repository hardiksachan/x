// Package xlog provides logging functions.
package xlog

// Message is a log message.
type Message struct {
	Title   string
	Details string
	Data    map[string]string
}

// Log is a logger.
type Log interface {
	Info(msg Message)
	Debug(msg Message)
	Warn(msg Message)
	Error(msg Message)
}
