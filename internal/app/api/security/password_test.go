package security

import (
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

func TestPlainParser_Parse(t *testing.T) {
	parser, found := parsers[passwordTypePlain]

	if !found {
		t.Error("parser not initialized")
	}

	t.Run("success", func(t *testing.T) {
		h, err := parser.Parse("password", 10)

		if err != nil {
			t.Error(err)
		}

		if err = bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid cost", func(t *testing.T) {
		h, err := parser.Parse("password", 999)

		if err == nil {
			t.Error("expected error, but got nil")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})
}

func TestEnvParser_Parse(t *testing.T) {
	parser, found := parsers[passwordTypeEnv]

	if !found {
		t.Error("parser not initialized")
	}

	if err := os.Setenv("TEST_PASSWORD", "password"); err != nil {
		t.Error(err)
	}

	t.Run("success", func(t *testing.T) {
		h, err := parser.Parse("TEST_PASSWORD", 10)

		if err != nil {
			t.Error(err)
		}

		if err = bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid cost", func(t *testing.T) {
		h, err := parser.Parse("password", 999)

		if err == nil {
			t.Error("expected error, but got nil")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})

	if err := os.Unsetenv("TEST_PASSWORD"); err != nil {
		t.Error(err)
	}
}

func TestEncryptParser_Parse(t *testing.T) {
	parser, found := parsers[passwordTypeEncrypt]

	if !found {
		t.Error("parser not initialized")
	}

	t.Run("success", func(t *testing.T) {
		h, err := parser.Parse("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC", 10)

		if err != nil {
			t.Error(err)
		}

		if err = bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid cost", func(t *testing.T) {
		h, err := parser.Parse("$2a$12$Yo2LBZZxseAXEDJYFDzlru.E3.PJeOm4HxEqolnNblXFH2vVt7crC", 10)

		if err == nil {
			t.Error("expected error, but got nil")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})
}

func TestParseHashedPassword(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		h, ok := ParseHashedPassword("", 10)

		if ok {
			t.Error("unexpected result")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})

	t.Run("illegal value", func(t *testing.T) {
		h, ok := ParseHashedPassword("test", 10)

		if ok {
			t.Error("unexpected result")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})

	t.Run("empty password hash", func(t *testing.T) {
		h, ok := ParseHashedPassword("!test:", 10)

		if ok {
			t.Error("unexpected result")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})

	t.Run("empty type", func(t *testing.T) {
		h, ok := ParseHashedPassword("!:test", 10)

		if ok {
			t.Error("unexpected result")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})

	t.Run("illegal type", func(t *testing.T) {
		h, ok := ParseHashedPassword("test:test", 10)

		if ok {
			t.Error("unexpected result")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})

	t.Run("password type / plain", func(t *testing.T) {
		h, ok := ParseHashedPassword("!plain:password", 10)

		if !ok {
			t.Error("unexpected result")
		}

		if err := bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("password type / env", func(t *testing.T) {
		if err := os.Setenv("TEST_PASSWORD", "password"); err != nil {
			t.Error(err)
		}

		h, ok := ParseHashedPassword("!env:TEST_PASSWORD", 10)

		if !ok {
			t.Error("unexpected result")
		}

		if err := bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}

		if err := os.Unsetenv("TEST_PASSWORD"); err != nil {
			t.Error(err)
		}
	})

	t.Run("password type / encrypt", func(t *testing.T) {
		h, ok := ParseHashedPassword("!encrypt:$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC", 10)

		if !ok {
			t.Error("unexpected result")
		}

		if err := bcrypt.CompareHashAndPassword(h, []byte("password")); err != nil {
			t.Error(err)
		}
	})

	t.Run("password type / not supported", func(t *testing.T) {
		h, ok := ParseHashedPassword("!test:test", 10)

		if ok {
			t.Error("unexpected result")
		}

		if h != nil {
			t.Errorf("expected nil hash, but got '%v'", h)
		}
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ok := VerifyPassword([]byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC"), []byte("password"))

		if !ok {
			t.Error("unexpected result")
		}
	})

	t.Run("failure", func(t *testing.T) {
		ok := VerifyPassword([]byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC"), []byte("test"))

		if ok {
			t.Error("unexpected result")
		}
	})
}
