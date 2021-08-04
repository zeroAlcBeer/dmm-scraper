package logger

import (
	"github.com/kataras/golog"
)

// Logger interface defines all the logging methods to be implemented
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Fatal(args ...interface{})

	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// New returns a new instance of Logger
func New() Logger {
	return golog.New()
}
