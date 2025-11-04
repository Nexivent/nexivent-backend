package logging

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Creates a new logger instance with the specified source, server type, log format and log level.
func NewLogger(source, serverType, logFormat Format, logLevel int) *LoggerWrapper {
	logger := logrus.New()
	logger.SetReportCaller(false)

	switch logLevel {
	case 0:
		logger.Level = logrus.PanicLevel
	case 1:
		logger.Level = logrus.FatalLevel
	case 2:
		logger.Level = logrus.ErrorLevel
	case 3:
		logger.Level = logrus.WarnLevel
	case 4:
		logger.Level = logrus.InfoLevel
	case 5:
		logger.Level = logrus.DebugLevel
	case 6:
		logger.Level = logrus.TraceLevel
	}

	switch logFormat {
	case FormatJSON:
		logger.Formatter = &logrus.JSONFormatter{CallerPrettyfier: CallerPrettyfier}
	case FormatText:
		logger.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: CallerPrettyfier,
			DisableTimestamp: true,
			QuoteEmptyFields: true,
		}
	}

	return &LoggerWrapper{
		log: logger.WithFields(logrus.Fields{
			"source": source,
			"server": serverType,
		}),
	}
}

// This function is required when you want to introduce your custom format.
// In this case file and line look like this `file="engine.go:141`
// but f.File provides a full path along with the file name.
// So in `formatFilePath()` function just trimmet everything before the file name
// and added a line number in the end.
func CallerPrettyfier(f *runtime.Frame) (string, string) {
	return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
