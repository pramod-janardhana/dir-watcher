package zlog

type noLogger struct {
}

// NewNoLogger creates a instance of of nologger, which ignores all log messages
func NewNoLogger() Logger {
	return &noLogger{}
}

func (n *noLogger) Infof(format string, v ...interface{})  {}
func (n *noLogger) Errorf(format string, v ...interface{}) {}
func (n *noLogger) Debugf(format string, v ...interface{}) {}
func (n *noLogger) Fatalf(format string, v ...interface{}) {}
