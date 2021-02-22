package logger

import (
	"github.com/sirupsen/logrus"
)

// LogrusLogger ...
type LogrusLogger struct {
	*logrus.Logger
}
