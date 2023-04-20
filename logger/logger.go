package logger

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type Level logrus.Level

const (
	FatalLevel Level = Level(logrus.FatalLevel)
	ErrorLevel Level = Level(logrus.ErrorLevel)
	WarnLevel  Level = Level(logrus.WarnLevel)
	InfoLevel  Level = Level(logrus.InfoLevel)
	DebugLevel Level = Level(logrus.DebugLevel)
)

var (
	logger    *logrus.Logger
	AllLevels = []Level{FatalLevel, ErrorLevel, WarnLevel, InfoLevel, DebugLevel}
)

func init() {
	logger = logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		CallerPrettyfier: getCallerPrettyfier(false),
	})
	SetLevel(InfoLevel)
}

func getCallerPrettyfier(fullpath bool) func(f *runtime.Frame) (string, string) {
	// https://github.com/sirupsen/logrus/blob/v1.9.0/example_custom_caller_test.go
	// https://github.com/kubernetes/klog/blob/v2.90.1/klog.go#L644
	if fullpath {
		return func(f *runtime.Frame) (string, string) {
			_, file, line, ok := runtime.Caller(9)
			if !ok {
				file = "???"
				line = 1
			}
			return "", fmt.Sprintf("%s:%d", file, line)
		}
	}
	return func(f *runtime.Frame) (string, string) {
		_, file, line, ok := runtime.Caller(9)
		if !ok {
			file = "???"
			line = 1
		} else {
			slash := strings.LastIndex(file, "/")
			if slash >= 0 {
				file = file[slash+1:]
			}
		}
		return "", fmt.Sprintf("%s:%d", file, line)
	}
}

// setters & getters

func SetLevel(level Level) {
	logger.SetLevel(logrus.Level(level))
}

func GetLevel() Level {
	return Level(logger.GetLevel())
}

func SetOutput(output io.Writer) {
	logger.SetOutput(output)
}

func SetFullpath(fullpath bool) {
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		CallerPrettyfier: getCallerPrettyfier(fullpath),
	})
}

// log functions

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
