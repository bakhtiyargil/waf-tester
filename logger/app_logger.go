package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/fs"
	"os"
	"syscall"
	"waf-tester/config"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"error": zapcore.ErrorLevel,
}

type Logger interface {
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	ErrorR(args ...interface{}) error
}

type AppLogger struct {
	sugarLogger *zap.SugaredLogger
}

func (l *AppLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *AppLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *AppLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *AppLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *AppLogger) ErrorR(args ...interface{}) error {
	l.sugarLogger.Error(args)
	return args[len(args)-1].(error)
}

func (l *AppLogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		level = zapcore.InfoLevel
	}
	return level
}

func InitAppLogger(cfg *config.Config) *AppLogger {
	l := &AppLogger{}
	logLevel := l.getLoggerLevel(cfg)
	logWriter := zapcore.AddSync(os.Stderr)
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()

	/*
		:uber-zap ENOTTY error bypass (no fix!)
	*/
	var pathErr *fs.PathError
	if err := l.sugarLogger.Sync(); !(errors.Is(err, syscall.ENOTTY) || errors.As(err, &pathErr)) {
		l.sugarLogger.Error(err)
	}
	return l
}
