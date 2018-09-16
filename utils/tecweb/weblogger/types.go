package weblogger

type Logger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Warningln(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Printf(format string, v ...interface{})

	Fatal(args ...interface{})
	Fatalln(args ...interface{})
	Fatalf(format string, args ...interface{})

	Exit(args ...interface{})
	Exitln(args ...interface{})
	Exitf(format string, args ...interface{})

	Flush()
}
