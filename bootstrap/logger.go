package bootstrap

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/fs"
	"os"
	"syscall"
	"waf-tester/config"
	"waf-tester/logger"
)

var loggerLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"error": zapcore.ErrorLevel,
}

func InitAppLogger() logger.Logger {
	l := &logger.AppLogger{}
	logLevel := getLoggerLevel(App.Config)
	logWriter := zapcore.AddSync(os.Stderr)
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	lggr := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.SugarLogger = lggr.Sugar()

	/*
		:uber-zap ENOTTY error bypass (no fix!)
	*/
	var pathErr *fs.PathError
	if err := l.SugarLogger.Sync(); !(errors.Is(err, syscall.ENOTTY) || errors.As(err, &pathErr)) {
		l.SugarLogger.Error(err)
	}
	return l
}

func getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		level = zapcore.InfoLevel
	}
	return level
}
