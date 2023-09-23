package xlog

type Message struct {
	Title   string
	Details string
	Data    map[string]string
}

type Log interface {
	Info(msg Message)
	Debug(msg Message)
	Warn(msg Message)
	Error(msg Message)
}
