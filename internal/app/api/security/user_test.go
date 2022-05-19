package security

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	var (
		expectedName           = "test"
		expectedHashedPassword = []byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC")
		expectedScopes         = []Scope{ScopeReadServer}
	)

	user := NewUser(expectedName, expectedHashedPassword, expectedScopes)

	if user == nil {
		t.Error("expected user, but got nil")
	}

	if name := user.Name(); !strings.EqualFold(expectedName, name) {
		t.Errorf("expected user name '%v', but got '%v'", expectedName, name)
	}

	if hashedPassword := user.HashedPassword(); !bytes.Equal(expectedHashedPassword, hashedPassword) {
		t.Errorf("expected user hashed password '%v', but got '%v'", expectedHashedPassword, hashedPassword)
	}

	if scopes := user.Scopes(); !reflect.DeepEqual(expectedScopes, scopes) {
		t.Errorf("expected user scopes '%v', but got '%v'", expectedScopes, scopes)
	}
}

func TestNewUserProvider(t *testing.T) {
	provider := NewUserProvider()

	provider.Set("test", NewUser("test", []byte{}, []Scope{}))

	t.Run("user / exists", func(t *testing.T) {
		user := provider.Get("test")

		if user == nil {
			t.Error("expected user, but got nil")
		}
	})

	t.Run("user / unknown", func(t *testing.T) {
		user := provider.Get("unknown")

		if user != nil {
			t.Errorf("expected nil, but got user '%v'", user)
		}
	})
}
