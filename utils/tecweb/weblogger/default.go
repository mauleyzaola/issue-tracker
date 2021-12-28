package weblogger

import (
	"github.com/golang/glog"
)

type DefaultLogger struct{}

func NewDefaultLogger() Logger {
	t := &DefaultLogger{}
	return t
}

func (t *DefaultLogger) Error(args ...interface{}) {
	glog.Error(args...)
}

func (t *DefaultLogger) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}

func (t *DefaultLogger) Errorln(args ...interface{}) {
	glog.Errorln(args...)
}

func (t *DefaultLogger) Warning(args ...interface{}) {
	glog.Warning(args...)
}

func (t *DefaultLogger) Warningf(format string, args ...interface{}) {
	glog.Warningf(format, args...)
}

func (t *DefaultLogger) Warningln(args ...interface{}) {
	glog.Warningln(args...)
}

func (t *DefaultLogger) Info(args ...interface{}) {
	glog.Info(args...)
}

func (t *DefaultLogger) Infof(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

func (t *DefaultLogger) Infoln(args ...interface{}) {
	glog.Infoln(args...)
}

func (t *DefaultLogger) Printf(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

func (t *DefaultLogger) Fatal(args ...interface{}) {
	glog.Fatal(args...)
}

func (t *DefaultLogger) Fatalln(args ...interface{}) {
	glog.Fatalln(args...)
}

func (t *DefaultLogger) Fatalf(format string, args ...interface{}) {
	glog.Fatalf(format, args...)
}

func (t *DefaultLogger) Exit(args ...interface{}) {
	glog.Exit(args...)
}

func (t *DefaultLogger) Exitln(args ...interface{}) {
	glog.Exitln(args...)
}

func (t *DefaultLogger) Exitf(format string, args ...interface{}) {
	glog.Exitf(format, args...)
}

func (t *DefaultLogger) Flush() {
	glog.Flush()
}
