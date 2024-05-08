package loader

type LogAdapter interface {
	Error(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
}

type LogNoop struct{}

func (LogNoop) Error(msg string, keysAndValues ...interface{}) {}
func (LogNoop) Info(msg string, keysAndValues ...interface{})  {}
func (LogNoop) Debug(msg string, keysAndValues ...interface{}) {}
func (LogNoop) Warn(msg string, keysAndValues ...interface{})  {}

var _ LogAdapter = LogNoop{}
