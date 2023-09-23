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

func Info(msg Message) {
	currentLogger().Info(msg)
}

func InfoString(msg string) {
	currentLogger().Info(Message{
		Title:   msg,
		Details: "",
		Data:    nil,
	})
}

func Infof(format string, args ...interface{}) {
	InfoString(fmt.Sprintf(format, args...))
}

func Warn(msg Message) {
	currentLogger().Warn(msg)
}

func WarnString(msg string) {
	currentLogger().Warn(Message{
		Title:   msg,
		Details: "",
		Data:    nil,
	})
}

func Warnf(format string, args ...interface{}) {
	WarnString(fmt.Sprintf(format, args...))
}

func Error(msg Message) {
	currentLogger().Error(msg)
}

func ErrorString(msg string) {
	currentLogger().Error(Message{
		Title:   msg,
		Details: "",
		Data:    nil,
	})
}

func Errorf(format string, args ...interface{}) {
	ErrorString(fmt.Sprintf(format, args...))
}

func Debug(msg Message) {
	currentLogger().Debug(msg)
}

func DebugString(msg string) {
	currentLogger().Debug(Message{
		Title:   msg,
		Details: "",
		Data:    nil,
	})
}

func Debugf(format string, args ...interface{}) {
	DebugString(fmt.Sprintf(format, args...))
}
