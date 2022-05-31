package security_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

func TestPlainParser_Parse(t *testing.T) {
	parser := security.MustGetPasswordParser(security.PasswordTypePlain)

	t.Run("success", func(t *testing.T) {
		h, err := parser.Parse("password", 10)

		require.Nil(t, err)

		if err = bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid cost", func(t *testing.T) {
		h, err := parser.Parse("password", 999)

		require.NotNil(t, err)
		require.Nil(t, h)
	})
}

func TestEnvParser_Parse(t *testing.T) {
	parser := security.MustGetPasswordParser(security.PasswordTypeEnv)

	if err := os.Setenv("TEST_PASSWORD", "password"); err != nil {
		t.Fatal(err)
	}

	t.Run("success", func(t *testing.T) {
		h, err := parser.Parse("TEST_PASSWORD", 10)

		require.Nil(t, err)

		if err = bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid cost", func(t *testing.T) {
		h, err := parser.Parse("password", 999)

		require.NotNil(t, err)
		require.Nil(t, h)
	})

	if err := os.Unsetenv("TEST_PASSWORD"); err != nil {
		t.Fatal(err)
	}
}

func TestEncryptParser_Parse(t *testing.T) {
	parser := security.MustGetPasswordParser(security.PasswordTypeEncrypt)

	t.Run("success", func(t *testing.T) {
		h, err := parser.Parse("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC", 10)

		require.Nil(t, err)

		if err = bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid cost", func(t *testing.T) {
		h, err := parser.Parse("$2a$12$Yo2LBZZxseAXEDJYFDzlru.E3.PJeOm4HxEqolnNblXFH2vVt7crC", 10)

		require.NotNil(t, err)
		require.Nil(t, h)
	})
}

func TestParseHashedPassword(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		h, ok := security.ParseHashedPassword("", 10)

		require.False(t, ok)
		require.Nil(t, h)
	})

	t.Run("illegal value", func(t *testing.T) {
		h, ok := security.ParseHashedPassword("test", 10)

		require.False(t, ok)
		require.Nil(t, h)
	})

	t.Run("empty password hash", func(t *testing.T) {
		h, ok := security.ParseHashedPassword("!test:", 10)

		require.False(t, ok)
		require.Nil(t, h)
	})

	t.Run("empty type", func(t *testing.T) {
		h, ok := security.ParseHashedPassword("!:test", 10)

		require.False(t, ok)
		require.Nil(t, h)
	})

	t.Run("illegal type", func(t *testing.T) {
		h, ok := security.ParseHashedPassword("test:test", 10)

		require.False(t, ok)
		require.Nil(t, h)
	})

	t.Run("password type / plain", func(t *testing.T) {
		h, ok := security.ParseHashedPassword("!plain:password", 10)

		require.True(t, ok)

		if err := bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("password type / env", func(t *testing.T) {
		if err := os.Setenv("TEST_PASSWORD", "password"); err != nil {
			t.Fatal(err)
		}

		h, ok := security.ParseHashedPassword("!env:TEST_PASSWORD", 10)

		require.True(t, ok)

		if err := bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}

		if err := os.Unsetenv("TEST_PASSWORD"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("password type / encrypt", func(t *testing.T) {
		h, ok := security.ParseHashedPassword("!encrypt:$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC", 10)

		require.True(t, ok)

		if err := bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("password type / not supported", func(t *testing.T) {
		h, ok := security.ParseHashedPassword("!test:test", 10)

		require.False(t, ok)
		require.Nil(t, h)
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ok := security.VerifyPassword([]byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC"), []byte("password"))

		require.True(t, ok)
	})

	t.Run("failure", func(t *testing.T) {
		ok := security.VerifyPassword([]byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC"), []byte("test"))

		require.False(t, ok)
	})
}
