package weblogger

import (
	"testing"
)

func TestLoggerImplementation(t *testing.T) {
	var l Logger
	l = NewDefaultLogger()
	l.Infoln("all ok")
}
