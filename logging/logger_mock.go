package logging

type LoggerMock struct{}

func NewLoggerMock() *LoggerMock {
	return &LoggerMock{}
}

func (l *LoggerMock) Fatal(format ...interface{}) {
}

func (l *LoggerMock) Fatalf(format string, args ...interface{}) {
}

func (l *LoggerMock) Fatalln(args ...interface{}) {
}

func (l *LoggerMock) Debug(args ...interface{}) {
}

func (l *LoggerMock) Debugf(format string, args ...interface{}) {
}

func (l *LoggerMock) Debugln(args ...interface{}) {
}

func (l *LoggerMock) Error(args ...interface{}) {
}

func (l *LoggerMock) Errorf(format string, args ...interface{}) {
}

func (l *LoggerMock) Errorln(args ...interface{}) {
}

func (l *LoggerMock) Info(args ...interface{}) {
}

func (l *LoggerMock) Infof(format string, args ...interface{}) {
}

func (l *LoggerMock) Infoln(args ...interface{}) {
}

func (l *LoggerMock) Warn(args ...interface{}) {
}

func (l *LoggerMock) Warnf(format string, args ...interface{}) {
}

func (l *LoggerMock) Warnln(args ...interface{}) {
}

func (l *LoggerMock) Panic(args ...interface{}) {
}

func (l *LoggerMock) Panicf(format string, args ...interface{}) {
}

func (l *LoggerMock) Panicln(args ...interface{}) {
}

func (l *LoggerMock) WithFields(fields map[string]interface{}) Logger {
	return &LoggerMock{}
}

func (l *LoggerMock) WithField(key string, value interface{}) Logger {
	return &LoggerMock{}
}

func (l *LoggerMock) WithError(err error) Logger {
	return &LoggerMock{}
}
