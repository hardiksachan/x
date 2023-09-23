package xlog

type noopLogger struct{}

func (z *noopLogger) Info(_ Message) {}

func (z *noopLogger) Warn(_ Message) {}

func (z *noopLogger) Error(_ Message) {}

func (z *noopLogger) Debug(_ Message) {}

func newNoopLogger() Log {
	return &noopLogger{}
}
