package log

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"go.uber.org/zap"
)

func NewLogger(config *config.Config) *zap.Logger {
	logger, err := NewLoggerConfig(config).Build()
	if err != nil {
		panic(err)
	}

	return logger
}
