package xlog

import (
	"fmt"
	"sync"
)

var (
	once     sync.Once
	instance Log
)

func currentLogger() Log {
	once.Do(func() {
		// TODO: initialise JSON currentLogger

		pretty, err := newPrettyLogger()
		if err != nil {
			instance = newNoopLogger()
		} else {
			instance = pretty
		}
	})

	return instance
}

// Info logs an info message.
func Info(msg Message) {
	currentLogger().Info(msg)
}

// InfoString logs an info message.
func InfoString(msg string) {
	currentLogger().Info(Message{
		Title:   msg,
		Details: "",
		Data:    nil,
	})
}

// Infof logs an info message.
func Infof(format string, args ...interface{}) {
	InfoString(fmt.Sprintf(format, args...))
}

// Warn logs a warning message.
func Warn(msg Message) {
	currentLogger().Warn(msg)
}

// WarnString logs a warning message.
func WarnString(msg string) {
	currentLogger().Warn(Message{
		Title:   msg,
		Details: "",
		Data:    nil,
	})
}

// Warnf logs a warning message.
func Warnf(format string, args ...interface{}) {
	WarnString(fmt.Sprintf(format, args...))
}

// Error logs an error message.
func Error(msg Message) {
	currentLogger().Error(msg)
}

// ErrorString logs an error message.
func ErrorString(msg string) {
	currentLogger().Error(Message{
		Title:   msg,
		Details: "",
		Data:    nil,
	})
}

// Errorf logs an error message.
func Errorf(format string, args ...interface{}) {
	ErrorString(fmt.Sprintf(format, args...))
}

// Debug logs a debug message.
func Debug(msg Message) {
	currentLogger().Debug(msg)
}

// DebugString logs a debug message.
func DebugString(msg string) {
	currentLogger().Debug(Message{
		Title:   msg,
		Details: "",
		Data:    nil,
	})
}

// Debugf logs a debug message.
func Debugf(format string, args ...interface{}) {
	DebugString(fmt.Sprintf(format, args...))
}
