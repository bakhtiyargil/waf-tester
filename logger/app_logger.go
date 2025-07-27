package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	ErrorR(args ...interface{}) error
	FatalF(template string, args ...interface{})
}

type AppLogger struct {
	SugarLogger *zap.SugaredLogger
}

func (l *AppLogger) Info(args ...interface{}) {
	l.SugarLogger.Info(args...)
}

func (l *AppLogger) Infof(template string, args ...interface{}) {
	l.SugarLogger.Infof(template, args...)
}

func (l *AppLogger) Error(args ...interface{}) {
	l.SugarLogger.Error(args...)
}

func (l *AppLogger) Errorf(template string, args ...interface{}) {
	l.SugarLogger.Errorf(template, args...)
}

func (l *AppLogger) ErrorR(args ...interface{}) error {
	l.SugarLogger.Error(args)
	return args[len(args)-1].(error)
}

func (l *AppLogger) FatalF(template string, args ...interface{}) {
	l.SugarLogger.Fatalf(template, args)
}
