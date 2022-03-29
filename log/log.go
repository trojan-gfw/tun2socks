package log

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

// _defaultLevel is package default logging level.
var _defaultLevel = atomic.NewUint32(uint32(InfoLevel))

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func SetLevel(level Level) {
	_defaultLevel.Store(uint32(level))
}

func Debugf(format string, args ...any) {
	logf(DebugLevel, format, args...)
}

func Infof(format string, args ...any) {
	logf(InfoLevel, format, args...)
}

func Warnf(format string, args ...any) {
	logf(WarnLevel, format, args...)
}

func Errorf(format string, args ...any) {
	logf(ErrorLevel, format, args...)
}

func Fatalf(format string, args ...any) {
	logf(FatalLevel, format, args...)
}

func logf(level Level, format string, args ...any) {
	event := newEvent(level, format, args...)
	if uint32(event.Level) > _defaultLevel.Load() {
		return
	}

	switch level {
	case DebugLevel:
		logrus.WithTime(event.Time).WithField("fromapp", "igniter-golib-log-tun2socks").Debugln(event.Message)
	case InfoLevel:
		logrus.WithTime(event.Time).WithField("fromapp", "igniter-golib-log-tun2socks").Infoln(event.Message)
	case WarnLevel:
		logrus.WithTime(event.Time).WithField("fromapp", "igniter-golib-log-tun2socks").Warnln(event.Message)
	case ErrorLevel:
		logrus.WithTime(event.Time).WithField("fromapp", "igniter-golib-log-tun2socks").Errorln(event.Message)
	case FatalLevel:
		logrus.WithTime(event.Time).WithField("fromapp", "igniter-golib-log-tun2socks").Fatalln(event.Message)
	}
}
