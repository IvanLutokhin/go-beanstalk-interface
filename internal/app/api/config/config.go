package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"time"
)

const (
	EnvVariableKey = "BEANSTALK_INTERFACE_CONFIG"
	DefaultFile    = "configs/api.config.yaml"
)

type Config struct {
	Logger    LoggerConfig    `yaml:"logger"`
	Beanstalk BeanstalkConfig `yaml:"beanstalk"`
	Http      HttpConfig      `yaml:"http"`
}

type LoggerConfig struct {
	Level         zap.AtomicLevel        `yaml:"level"`
	Encoding      string                 `yaml:"encoding"`
	Encoder       EncoderConfig          `yaml:"encoder"`
	InitialFields map[string]interface{} `yaml:"initial_fields"`
}

type EncoderConfig struct {
	MessageKey       string                  `yaml:"message_key"`
	LevelKey         string                  `yaml:"level_key"`
	TimeKey          string                  `yaml:"time_key"`
	NameKey          string                  `yaml:"name_key"`
	CallerKey        string                  `yaml:"caller_key"`
	FunctionKey      string                  `yaml:"function_key"`
	StacktraceKey    string                  `yaml:"stacktrace_key"`
	LineEnding       string                  `yaml:"line_ending"`
	LevelEncoder     zapcore.LevelEncoder    `yaml:"level_encoder"`
	TimerEncoder     zapcore.TimeEncoder     `yaml:"time_encoder"`
	DurationEncoder  zapcore.DurationEncoder `yaml:"duration_encoder"`
	CallerEncoder    zapcore.CallerEncoder   `yaml:"caller_encoder"`
	NameEncoder      zapcore.NameEncoder     `yaml:"name_encoder"`
	ConsoleSeparator string                  `yaml:"console_separator"`
}

type BeanstalkConfig struct {
	Address string     `yaml:"address"`
	Pool    PoolConfig `yaml:"pool"`
}

type PoolConfig struct {
	Capacity int `yaml:"capacity"`
}

type HttpConfig struct {
	ListenAddresses string        `yaml:"listen_addresses"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
}

func New() *Config {
	file, ok := os.LookupEnv(EnvVariableKey)
	if !ok {
		file = DefaultFile
	}

	c, err := Load(file)
	if err != nil {
		panic(err)
	}

	return c
}

func Load(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return Unmarshal(f)
}

func Unmarshal(reader io.Reader) (*Config, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var c Config
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
