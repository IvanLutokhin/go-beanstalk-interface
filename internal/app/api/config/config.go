package config

import (
	"errors"
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"time"
)

type Config struct {
	Logger    LoggerConfig    `yaml:"logger"`
	Beanstalk BeanstalkConfig `yaml:"beanstalk"`
	Http      HttpConfig      `yaml:"http"`
	Security  SecurityConfig  `yaml:"security"`
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
	Capacity    int           `yaml:"capacity"`
	MaxAge      time.Duration `yaml:"max_age"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type HttpConfig struct {
	ListenAddresses string        `yaml:"listen_addresses"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
}

type SecurityConfig struct {
	Secret     string        `yaml:"secret"`
	TokenTTL   time.Duration `yaml:"token_ttl"`
	BCryptCost int           `yaml:"bcrypt_cost"`
	Users      []UserConfig  `yaml:"users"`
}

type UserConfig struct {
	Name     string   `yaml:"name"`
	Password string   `yaml:"password"`
	Scopes   []string `yaml:"scopes"`
}

func LoadOrDefault(name string) (*Config, error) {
	if c, err := Load(name); err == nil {
		return c, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else {
		return Default(), nil
	}
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

func Default() *Config {
	return &Config{
		Logger: LoggerConfig{
			Level:    zap.NewAtomicLevelAt(zapcore.DebugLevel),
			Encoding: "console",
			Encoder: EncoderConfig{
				MessageKey:       "message",
				LevelKey:         "level",
				TimeKey:          "timestamp",
				NameKey:          "logger",
				CallerKey:        "caller",
				FunctionKey:      "",
				StacktraceKey:    "stacktrace",
				LineEnding:       "\n",
				LevelEncoder:     zapcore.CapitalLevelEncoder,
				TimerEncoder:     zapcore.RFC3339TimeEncoder,
				DurationEncoder:  zapcore.StringDurationEncoder,
				CallerEncoder:    zapcore.ShortCallerEncoder,
				NameEncoder:      zapcore.FullNameEncoder,
				ConsoleSeparator: "\t",
			},
			InitialFields: make(map[string]interface{}),
		},
		Beanstalk: BeanstalkConfig{
			Address: env.MustGetString("BI_SERVER_ADDRESS", "127.0.0.1:11300"),
			Pool: PoolConfig{
				Capacity:    env.MustGetInt("BI_POOL_CAPACITY", 25),
				MaxAge:      0,
				IdleTimeout: 0,
			},
		},
		Http: HttpConfig{
			ListenAddresses: env.MustGetString("BI_LISTEN_ADDRESSES", ":9999"),
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     60 * time.Second,
		},
		Security: SecurityConfig{
			Secret:     env.MustGetString("BI_SECURITY_SECRET", "secret"),
			TokenTTL:   time.Hour,
			BCryptCost: env.MustGetInt("BI_SECURITY_BCRYPT_COST", 10),
			Users: []UserConfig{
				{
					Name:     env.MustGetString("BI_ROOT_USER", "admin"),
					Password: env.MustGetString("BI_ROOT_PASSWORD", "!plain:admin"),
					Scopes: []string{
						"read:server",
						"read:tubes",
						"read:jobs",
						"write:jobs",
					},
				},
			},
		},
	}
}
