package config_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestLoadOrDefault(t *testing.T) {
	t.Run("load", func(t *testing.T) {
		configPath := "../../../../test/testdata/api.config.yaml"

		if _, err := os.Stat(configPath); err != nil {
			t.Fatal(err)
		}

		c, err := config.LoadOrDefault(configPath)

		require.Nil(t, err)
		require.NotNil(t, c, "Expected file config")

		// Logger
		require.Equal(t, "debug", c.Logger.Level.String(), "logger.level")
		require.Equal(t, "console", c.Logger.Encoding, "logger.encoding")
		require.Equal(t, "message", c.Logger.Encoder.MessageKey, "logger.encoder.message_key")
		require.Equal(t, "level", c.Logger.Encoder.LevelKey, "logger.encoder.level_key")
		require.Equal(t, "timestamp", c.Logger.Encoder.TimeKey, "logger.encoder.time_key")
		require.Equal(t, "logger", c.Logger.Encoder.NameKey, "logger.encoder.name_key")
		require.Equal(t, "caller", c.Logger.Encoder.CallerKey, "logger.encoder.caller_key")
		require.Equal(t, "", c.Logger.Encoder.FunctionKey, "logger.encoder.function_key")
		require.Equal(t, "stacktrace", c.Logger.Encoder.StacktraceKey, "logger.encoder.stacktrace_key")
		require.Equal(t, "\n", c.Logger.Encoder.LineEnding, "logger.encoder.line_ending")
		require.Equal(t, "\t", c.Logger.Encoder.ConsoleSeparator, "logger.encoder.console_separator")
		require.Len(t, c.Logger.InitialFields, 1, "logger.initial_fields")
		require.Equal(t, "value", c.Logger.InitialFields["key"].(string), "logger.initial_fields[key]")

		// Beanstalk
		require.Equal(t, "127.0.0.1:11300", c.Beanstalk.Address, "beanstalk.address")
		require.Equal(t, 3, c.Beanstalk.Pool.Capacity, "beanstalk.pool.capacity")

		// Http
		require.Equal(t, ":9999", c.Http.ListenAddresses, "http.listen_addresses")
		require.Equal(t, 30*time.Second, c.Http.ReadTimeout, "http.read_timeout")
		require.Equal(t, 30*time.Second, c.Http.WriteTimeout, "http.write_timeout")
		require.Equal(t, 60*time.Second, c.Http.IdleTimeout, "http.idle_timeout")
		require.ElementsMatch(t, []string{"*"}, c.Http.Cors.AllowedOrigins, "http.cors.allow_origins")
		require.ElementsMatch(t, []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"}, c.Http.Cors.AllowedMethods, "http.cors.allow_methods")
		require.ElementsMatch(t, []string{"Accept", "Authorization", "Content-Type", "Origin", "X-Requested-With"}, c.Http.Cors.AllowedHeaders, "http.cors.allow_headers")
		require.True(t, c.Http.Cors.AllowCredentials, "http.cors.allow_credentials")

		// Security
		require.Equal(t, 10, c.Security.BCryptCost, "security.bcrypt_cost")
		require.Len(t, c.Security.Users, 1, "security.users")
		require.Equal(t, "admin", c.Security.Users[0].Name, "security.users[0].name")
		require.Equal(t, "!plain:admin", c.Security.Users[0].Password, "security.users[0].password")
		require.ElementsMatch(t, []string{"read:server", "read:tubes", "read:jobs", "write:jobs"}, c.Security.Users[0].Scopes, "security.users[0].scopes")
	})

	t.Run("default", func(t *testing.T) {
		c, err := config.LoadOrDefault("")

		require.Nil(t, err)
		require.NotNil(t, c, "Expected default config")

		// Logger
		require.Equal(t, "debug", c.Logger.Level.String(), "logger.level")
		require.Equal(t, "console", c.Logger.Encoding, "logger.encoding")
		require.Equal(t, "message", c.Logger.Encoder.MessageKey, "logger.encoder.message_key")
		require.Equal(t, "level", c.Logger.Encoder.LevelKey, "logger.encoder.level_key")
		require.Equal(t, "timestamp", c.Logger.Encoder.TimeKey, "logger.encoder.time_key")
		require.Equal(t, "logger", c.Logger.Encoder.NameKey, "logger.encoder.name_key")
		require.Equal(t, "caller", c.Logger.Encoder.CallerKey, "logger.encoder.caller_key")
		require.Equal(t, "", c.Logger.Encoder.FunctionKey, "logger.encoder.function_key")
		require.Equal(t, "stacktrace", c.Logger.Encoder.StacktraceKey, "logger.encoder.stacktrace_key")
		require.Equal(t, "\n", c.Logger.Encoder.LineEnding, "logger.encoder.line_ending")
		require.Equal(t, "\t", c.Logger.Encoder.ConsoleSeparator, "logger.encoder.console_separator")
		require.Len(t, c.Logger.InitialFields, 0, "logger.initial_fields")

		// Beanstalk
		require.Equal(t, "127.0.0.1:11300", c.Beanstalk.Address, "beanstalk.address")
		require.Equal(t, 3, c.Beanstalk.Pool.Capacity, "beanstalk.pool.capacity")

		// Http
		require.Equal(t, ":9999", c.Http.ListenAddresses, "http.listen_addresses")
		require.Equal(t, 30*time.Second, c.Http.ReadTimeout, "http.read_timeout")
		require.Equal(t, 30*time.Second, c.Http.WriteTimeout, "http.write_timeout")
		require.Equal(t, 60*time.Second, c.Http.IdleTimeout, "http.idle_timeout")
		require.ElementsMatch(t, []string{"*"}, c.Http.Cors.AllowedOrigins, "http.cors.allow_origins")
		require.ElementsMatch(t, []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"}, c.Http.Cors.AllowedMethods, "http.cors.allow_methods")
		require.ElementsMatch(t, []string{"Accept", "Authorization", "Content-Type", "Origin", "X-Requested-With"}, c.Http.Cors.AllowedHeaders, "http.cors.allow_headers")
		require.True(t, c.Http.Cors.AllowCredentials, "http.cors.allow_credentials")

		// Security
		require.Equal(t, 10, c.Security.BCryptCost, "security.bcrypt_cost")
		require.Len(t, c.Security.Users, 1, "security.users")
		require.Equal(t, "admin", c.Security.Users[0].Name, "security.users[0].name")
		require.Equal(t, "!plain:admin", c.Security.Users[0].Password, "security.users[0].password")
		require.ElementsMatch(t, []string{"read:server", "read:tubes", "read:jobs", "write:jobs"}, c.Security.Users[0].Scopes, "security.users[0].scopes")
	})
}
