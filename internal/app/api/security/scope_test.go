package security

import (
	"reflect"
	"testing"
)

func TestIsAvailableScope(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := IsAvailableScope("read:server")

		if !r {
			t.Error("unexpected result")
		}
	})

	t.Run("failure", func(t *testing.T) {
		r := IsAvailableScope("test")

		if r {
			t.Error("unexpected result")
		}
	})
}

func TestParseScopes(t *testing.T) {
	expectedScopes := []Scope{ScopeReadServer, ScopeReadTubes, ScopeReadJobs}

	scopes := ParseScopes([]string{"read:server", "read:server", "READ:TUBES", "READ:jobs", "test"})

	if !reflect.DeepEqual(expectedScopes, scopes) {
		t.Errorf("expected scopes '%v', but got '%v'", expectedScopes, scopes)
	}
}

func TestVerifyScopes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		scopes := VerifyScopes([]Scope{ScopeReadServer}, []Scope{ScopeReadServer})

		if len(scopes) != 0 {
			t.Errorf("expected empty scope slice, but got '%v'", scopes)
		}
	})

	t.Run("failure", func(t *testing.T) {
		expectedScopes := []Scope{ScopeWriteJobs}

		scopes := VerifyScopes([]Scope{ScopeReadServer, ScopeReadTubes, ScopeReadJobs}, []Scope{ScopeWriteJobs, "test"})

		if !reflect.DeepEqual(expectedScopes, scopes) {
			t.Errorf("expected scope slice '%v', but got '%v'", expectedScopes, scopes)
		}
	})
}
