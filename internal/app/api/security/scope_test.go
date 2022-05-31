package security_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"reflect"
	"testing"
)

func TestIsAvailableScope(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := security.IsAvailableScope("read:server")

		if !r {
			t.Error("unexpected result")
		}
	})

	t.Run("failure", func(t *testing.T) {
		r := security.IsAvailableScope("test")

		if r {
			t.Error("unexpected result")
		}
	})
}

func TestParseScopes(t *testing.T) {
	expectedScopes := []security.Scope{security.ScopeReadServer, security.ScopeReadTubes, security.ScopeReadJobs}

	scopes := security.ParseScopes([]string{"read:server", "read:server", "READ:TUBES", "READ:jobs", "test"})

	if !reflect.DeepEqual(expectedScopes, scopes) {
		t.Errorf("expected scopes '%v', but got '%v'", expectedScopes, scopes)
	}
}

func TestVerifyScopes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		scopes := security.VerifyScopes([]security.Scope{security.ScopeReadServer}, []security.Scope{security.ScopeReadServer})

		if len(scopes) != 0 {
			t.Errorf("expected empty scope slice, but got '%v'", scopes)
		}
	})

	t.Run("failure", func(t *testing.T) {
		expectedScopes := []security.Scope{security.ScopeWriteJobs}

		scopes := security.VerifyScopes([]security.Scope{security.ScopeReadServer, security.ScopeReadTubes, security.ScopeReadJobs}, []security.Scope{security.ScopeWriteJobs, "test"})

		if !reflect.DeepEqual(expectedScopes, scopes) {
			t.Errorf("expected scope slice '%v', but got '%v'", expectedScopes, scopes)
		}
	})
}
