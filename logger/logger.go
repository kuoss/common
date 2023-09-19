package logger

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	logger     *logrus.Logger
	AllLevels      = []Level{PanicLevel, FatalLevel, ErrorLevel, WarnLevel, InfoLevel, DebugLevel, TraceLevel}
	callerSkip int = 10 // 10 for prod(default), maybe 9 for goroutine or test code
)

func init() {
	logger = logrus.New()
	logger.SetReportCaller(true)
	SetLevel(InfoLevel)
	SetFullpath(false)
}

// setters & getters...

func SetOutput(output io.Writer) {
	logger.SetOutput(output)
}

func SetLevel(level Level) {
	logger.SetLevel(logrus.Level(level))
}

func GetLevel() Level {
	return Level(logger.GetLevel())
}

func SetFullpath(fullpath bool) {
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		CallerPrettyfier: getCallerPrettyfier(fullpath),
	})
}

func SetCallerSkip(skip int) {
	callerSkip = skip
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
		_, file, line, ok := runtime.Caller(callerSkip)
		if !ok {
			file = "???"
			line = 1
		} else {
			if slash := strings.LastIndex(file, "/"); slash >= 0 {
				file = file[slash+1:]
			}
		}
		return "", fmt.Sprintf("%s:%d", file, line)
	}
}

// log functions...

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
