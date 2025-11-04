package logging

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Fatal(format ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
	WithFields(fields map[string]interface{}) Logger
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
}

type LoggerWrapper struct {
	log logrus.FieldLogger
}

func (l *LoggerWrapper) Fatal(format ...interface{}) {
	l.log.Fatal(format...)
}

func (l *LoggerWrapper) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *LoggerWrapper) Fatalln(args ...interface{}) {
	l.log.Fatalln(args...)
}

func (l *LoggerWrapper) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *LoggerWrapper) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *LoggerWrapper) Debugln(args ...interface{}) {
	l.log.Debugln(args...)
}

func (l *LoggerWrapper) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *LoggerWrapper) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *LoggerWrapper) Errorln(args ...interface{}) {
	l.log.Errorln(args...)
}

func (l *LoggerWrapper) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *LoggerWrapper) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *LoggerWrapper) Infoln(args ...interface{}) {
	l.log.Infoln(args...)
}

func (l *LoggerWrapper) Warn(args ...interface{}) {
	l.log.Warn(args...)
}

func (l *LoggerWrapper) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *LoggerWrapper) Warnln(args ...interface{}) {
	l.log.Warnln(args...)
}

func (l *LoggerWrapper) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

func (l *LoggerWrapper) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

func (l *LoggerWrapper) Panicln(args ...interface{}) {
	l.log.Panicln(args...)
}

func (l *LoggerWrapper) WithFields(fields map[string]interface{}) Logger {
	return &LoggerWrapper{log: l.log.WithFields(fields)}
}

func (l *LoggerWrapper) WithField(key string, value interface{}) Logger {
	return &LoggerWrapper{log: l.log.WithField(key, value)}
}

func (l *LoggerWrapper) WithError(err error) Logger {
	return &LoggerWrapper{log: l.log.WithError(err)}
}
