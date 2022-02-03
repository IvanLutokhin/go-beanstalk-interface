package config

import (
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	s := "logger:\n" +
		"  level: 'debug'\n" +
		"  encoding: 'console'\n" +
		"  encoder:\n" +
		"    message_key: 'message'\n" +
		"    level_key: 'level'\n" +
		"    time_key: 'timestamp'\n" +
		"    name_key: 'logger'\n" +
		"    caller_key: 'caller'\n" +
		"    function_key: ''\n" +
		"    stacktrace_key: 'stacktrace'\n" +
		"    line_ending: '\\n'\n" +
		"    level_encoder: 'capital'\n" +
		"    time_encoder: 'rfc3339'\n" +
		"    duration_encoder: 'string'\n" +
		"    caller_encoder: 'short'\n" +
		"    name_encoder: 'full'\n" +
		"    console_separator: '\\t'\n" +
		"  initial_fields:\n" +
		"    key: 'value'\n\n" +
		"beanstalk:\n" +
		"  address: '127.0.0.1:11300'\n" +
		"  pool:\n" +
		"    capacity: 3\n\n" +
		"http:\n" +
		"  listen_addresses: ':9999'\n" +
		"  read_timeout: 30\n" +
		"  write_timeout: 30\n" +
		"  idle_timeout: 60\n\n"

	c, err := Unmarshal(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}

	if !strings.EqualFold("debug", c.Logger.Level.String()) {
		t.Errorf("logger.level: expected 'debug', but got '%s'", c.Logger.Level.String())
	}

	if !strings.EqualFold("console", c.Logger.Encoding) {
		t.Errorf("logger.encoding: expected 'console', but got '%s'", c.Logger.Encoding)
	}

	if !strings.EqualFold("message", c.Logger.Encoder.MessageKey) {
		t.Errorf("logger.encoder.message_key: expected 'message', but got '%s'", c.Logger.Encoder.MessageKey)
	}

	if !strings.EqualFold("level", c.Logger.Encoder.LevelKey) {
		t.Errorf("logger.encoder.level_key: expected 'level', but got '%s'", c.Logger.Encoder.LevelKey)
	}

	if !strings.EqualFold("timestamp", c.Logger.Encoder.TimeKey) {
		t.Errorf("logger.encoder.time_key: expected 'timestamp', but got '%s'", c.Logger.Encoder.TimeKey)
	}

	if !strings.EqualFold("logger", c.Logger.Encoder.NameKey) {
		t.Errorf("logger.encoder.name_key: expected 'logger', but got '%s'", c.Logger.Encoder.NameKey)
	}

	if !strings.EqualFold("caller", c.Logger.Encoder.CallerKey) {
		t.Errorf("logger.encoder.caller_key: expected 'caller', but got '%s'", c.Logger.Encoder.CallerKey)
	}

	if !strings.EqualFold("", c.Logger.Encoder.FunctionKey) {
		t.Errorf("logger.encoder.function_key: expected '', but got '%s'", c.Logger.Encoder.FunctionKey)
	}

	if !strings.EqualFold("stacktrace", c.Logger.Encoder.StacktraceKey) {
		t.Errorf("logger.encoder.stacktrace_key: expected 'console', but got '%s'", c.Logger.Encoder.StacktraceKey)
	}

	if !strings.EqualFold("\\n", c.Logger.Encoder.LineEnding) {
		t.Errorf("logger.encoder.line_ending: expected '\n', but got '%s'", c.Logger.Encoder.LineEnding)
	}

	if !strings.EqualFold("\\t", c.Logger.Encoder.ConsoleSeparator) {
		t.Errorf("logger.encoder.console_separator: expected '\\t', but got '%s'", c.Logger.Encoder.ConsoleSeparator)
	}

	if len(c.Logger.InitialFields) != 1 {
		t.Errorf("logger.initial_fields: expected len 1, but got %d", len(c.Logger.InitialFields))
	}

	if !strings.EqualFold("value", c.Logger.InitialFields["key"].(string)) {
		t.Errorf("logger.initial_fields[key]: expected 'value', but got '%s'", c.Logger.InitialFields["key"].(string))
	}

	if !strings.EqualFold("127.0.0.1:11300", c.Beanstalk.Address) {
		t.Errorf("beanstalk.address: expected 'console', but got '%s'", c.Beanstalk.Address)
	}

	if c.Beanstalk.Pool.Capacity != 3 {
		t.Errorf("beanstalk.pook.capacity: expected '3', but got '%d'", c.Beanstalk.Pool.Capacity)
	}

	if !strings.EqualFold(":9999", c.Http.ListenAddresses) {
		t.Errorf("http.listen_addresses: expected ':9999', but got '%s'", c.Http.ListenAddresses)
	}

	if c.Http.ReadTimeout != 30 {
		t.Errorf("http.read_timeout: expected '30', but got '%d'", c.Http.ReadTimeout)
	}

	if c.Http.WriteTimeout != 30 {
		t.Errorf("http.write_timeout: expected '30', but got '%d'", c.Http.WriteTimeout)
	}

	if c.Http.IdleTimeout != 60 {
		t.Errorf("http.idle_timeout: expected '60', but got '%d'", c.Http.IdleTimeout)
	}
}
