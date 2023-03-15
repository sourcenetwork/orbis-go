package logger

type Logger interface {
	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
	Fatalf(fmt string, args ...interface{})

	// Set the log level to one of "debug", "info", "warn", "error", "fatal", "panic".
	// If the level is not recognized, the default level is "debug".
	SetLevel(lv string)

	Sync() error
}
