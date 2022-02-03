package log

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLoggerConfig(config *config.Config) zap.Config {
	return zap.Config{
		Level:       config.Logger.Level,
		Development: false,
		Encoding:    config.Logger.Encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:       config.Logger.Encoder.MessageKey,
			LevelKey:         config.Logger.Encoder.LevelKey,
			TimeKey:          config.Logger.Encoder.TimeKey,
			NameKey:          config.Logger.Encoder.NameKey,
			CallerKey:        config.Logger.Encoder.CallerKey,
			FunctionKey:      config.Logger.Encoder.FunctionKey,
			StacktraceKey:    config.Logger.Encoder.StacktraceKey,
			LineEnding:       config.Logger.Encoder.LineEnding,
			EncodeLevel:      config.Logger.Encoder.LevelEncoder,
			EncodeTime:       config.Logger.Encoder.TimerEncoder,
			EncodeDuration:   config.Logger.Encoder.DurationEncoder,
			EncodeCaller:     config.Logger.Encoder.CallerEncoder,
			EncodeName:       config.Logger.Encoder.NameEncoder,
			ConsoleSeparator: config.Logger.Encoder.ConsoleSeparator,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    config.Logger.InitialFields,
	}
}
