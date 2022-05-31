package env_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/env"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetString(t *testing.T) {
	if err := os.Setenv("TEST_KEY", "test"); err != nil {
		t.Fatal(err)
	}

	t.Run("key exists", func(t *testing.T) {
		v, err := env.GetString("TEST_KEY")

		require.Nil(t, err)
		require.Equal(t, "test", v)
	})

	t.Run("key not exists", func(t *testing.T) {
		v, err := env.GetString("TEST_UNDEFINED_KEY")

		require.Equal(t, env.ErrVarNotExists, err)
		require.Empty(t, v)
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

		require.Equal(t, "test", v)
	})

	t.Run("key not exists", func(t *testing.T) {
		v := env.MustGetString("TEST_UNDEFINED_KEY", "default")

		require.Equal(t, "default", v)
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

		require.Nil(t, err)
		require.Equal(t, 1, v)
	})

	t.Run("key not exists", func(t *testing.T) {
		v, err := env.GetInt("TEST_UNDEFINED_KEY")

		require.Equal(t, env.ErrVarNotExists, err)
		require.Equal(t, 0, v)
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

		require.Equal(t, 1, v)
	})

	t.Run("key not exists", func(t *testing.T) {
		v := env.MustGetInt("TEST_UNDEFINED_KEY", 999)

		require.Equal(t, 999, v)
	})

	if err := os.Unsetenv("TEST_KEY"); err != nil {
		t.Fatal(err)
	}
}
