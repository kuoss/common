package logger

import (
	"github.com/sirupsen/logrus"
)

type Level logrus.Level

const (
	PanicLevel Level = Level(logrus.PanicLevel)
	FatalLevel Level = Level(logrus.FatalLevel)
	ErrorLevel Level = Level(logrus.ErrorLevel)
	WarnLevel  Level = Level(logrus.WarnLevel)
	InfoLevel  Level = Level(logrus.InfoLevel)
	DebugLevel Level = Level(logrus.DebugLevel)
	TraceLevel Level = Level(logrus.TraceLevel)
)

func (level Level) String() string {
	return logrus.Level(level).String()
}

func ParseLevel(lvl string) (Level, error) {
	level, err := logrus.ParseLevel(lvl)
	return Level(level), err
}
