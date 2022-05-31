package env_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/env"
	"os"
	"strings"
	"testing"
)

func TestGetString(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "test"); err != nil {
		t.Fatal(err)
	}

	t.Run("key exists", func(t *testing.T) {
		v, err := env.GetString("TEST_KEY")

		if err != nil {
			t.Errorf("expected nil, but got error '%v'", err)
		}

		if !strings.EqualFold("test", v) {
			t.Errorf("expected value 'test', but got '%v'", v)
		}
	})

	t.Run("key not exists", func(t *testing.T) {
		v, err := env.GetString("TEST_UNDEFINED_KEY")

		if err == nil {
			t.Errorf("expected error '%v', but got '%v'", env.ErrVarNotExists, err)
		}

		if !strings.EqualFold("", v) {
			t.Errorf("expected empty value, but got '%v'", v)
		}
	})

	if err := os.Unsetenv("TEST_KEY"); err != nil {
		t.Fatal(err)
	}
}

func TestMustGetString(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "test"); err != nil {
		t.Fatal(err)
	}

	t.Run("key exists", func(t *testing.T) {
		v := env.MustGetString("TEST_KEY", "default")

		if !strings.EqualFold("test", v) {
			t.Errorf("expected value 'test', but got '%v'", v)
		}
	})

	t.Run("key not exists", func(t *testing.T) {
		v := env.MustGetString("TEST_UNDEFINED_KEY", "default")

		if !strings.EqualFold("default", v) {
			t.Errorf("expected value 'default', but got '%v'", v)
		}
	})

	if err := os.Unsetenv("TEST_KEY"); err != nil {
		t.Fatal(err)
	}
}

func TestGetInt(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "1"); err != nil {
		t.Fatal(err)
	}

	t.Run("key exists", func(t *testing.T) {
		v, err := env.GetInt("TEST_KEY")

		if err != nil {
			t.Errorf("expected nil, but got error '%v'", err)
		}

		if v != 1 {
			t.Errorf("expected value '1', but got '%v'", v)
		}
	})

	t.Run("key not exists", func(t *testing.T) {
		v, err := env.GetInt("TEST_UNDEFINED_KEY")

		if err == nil {
			t.Errorf("expected error '%v', but got '%v'", env.ErrVarNotExists, err)
		}

		if v != 0 {
			t.Errorf("expected value '0', but got '%v'", v)
		}
	})

	if err := os.Unsetenv("TEST_KEY"); err != nil {
		t.Fatal(err)
	}
}

func TestMustGetInt(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "1"); err != nil {
		t.Fatal(err)
	}

	t.Run("key exists", func(t *testing.T) {
		v := env.MustGetInt("TEST_KEY", 999)

		if v != 1 {
			t.Errorf("expected value '1', but got '%v'", v)
		}
	})

	t.Run("key not exists", func(t *testing.T) {
		v := env.MustGetInt("TEST_UNDEFINED_KEY", 999)

		if v != 999 {
			t.Errorf("expected value '999', but got '%v'", v)
		}
	})

	if err := os.Unsetenv("TEST_KEY"); err != nil {
		t.Fatal(err)
	}
}
