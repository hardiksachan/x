package xlog

import (
	"go.uber.org/zap"
)

type prettyLogger struct {
	sugar *zap.SugaredLogger
}

func (z *prettyLogger) Info(msg Message) {
	z.sugar.Info(msg)
}

func (z *prettyLogger) Warn(msg Message) {
	z.sugar.Warn(msg)
}

func (z *prettyLogger) Error(msg Message) {
	z.sugar.Error(msg)
}

func (z *prettyLogger) Debug(msg Message) {
	z.sugar.Debug(msg)
}

func newPrettyLogger() (Log, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()
	return &prettyLogger{sugar}, nil
}
