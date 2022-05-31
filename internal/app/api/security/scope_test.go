package security_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsAvailableScope(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ok := security.IsAvailableScope("read:server")

		require.True(t, ok)
	})

	t.Run("failure", func(t *testing.T) {
		ok := security.IsAvailableScope("test")

		require.False(t, ok)
	})
}

func TestParseScopes(t *testing.T) {
	scopes := security.ParseScopes([]string{"read:server", "read:server", "READ:TUBES", "READ:jobs", "test"})

	require.ElementsMatch(t, []security.Scope{security.ScopeReadServer, security.ScopeReadTubes, security.ScopeReadJobs}, scopes)
}

func TestVerifyScopes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		scopes := security.VerifyScopes([]security.Scope{security.ScopeReadServer}, []security.Scope{security.ScopeReadServer})

		require.Len(t, scopes, 0)
	})

	t.Run("failure", func(t *testing.T) {
		scopes := security.VerifyScopes([]security.Scope{security.ScopeReadServer, security.ScopeReadTubes, security.ScopeReadJobs}, []security.Scope{security.ScopeWriteJobs, "test"})

		require.ElementsMatch(t, []security.Scope{security.ScopeWriteJobs}, scopes)
	})
}
