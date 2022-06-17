package beanstalk

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"go.uber.org/zap"
)

type LoggerAdapter struct {
	logger *zap.Logger
}

func NewLoggerAdapter(logger *zap.Logger) *LoggerAdapter {
	return &LoggerAdapter{
		logger: logger.Named("beanstalk"),
	}
}

func (l *LoggerAdapter) Log(level beanstalk.LogLevel, msg string, args map[string]interface{}) {
	fields := make([]zap.Field, 0, len(args))
	for k, v := range args {
		fields = append(fields, zap.Any(k, v))
	}

	switch level {
	case beanstalk.DebugLogLevel:
		l.logger.Debug(msg, fields...)
	case beanstalk.InfoLogLevel:
		l.logger.Info(msg, fields...)
	case beanstalk.WarningLogLevel:
		l.logger.Warn(msg, fields...)
	case beanstalk.ErrorLogLevel:
		l.logger.Error(msg, fields...)
	case beanstalk.PanicLogLevel:
		l.logger.Panic(msg, fields...)
	case beanstalk.FatalLogLevel:
		l.logger.Fatal(msg, fields...)
	}
}
