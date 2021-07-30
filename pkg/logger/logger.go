package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger interface defines all the logging methods to be implemented
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// New returns a new instance of Logger
func New() Logger {
	return &LogrusLogger{
		&logrus.Logger{
			Out:       os.Stdout,
			Formatter: new(logrus.JSONFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
	}
}
