package security_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewUser(t *testing.T) {
	var (
		expectedName           = "test"
		expectedHashedPassword = []byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC")
		expectedScopes         = []security.Scope{security.ScopeReadServer}
	)

	user := security.NewUser(expectedName, expectedHashedPassword, expectedScopes)

	require.NotNil(t, user)
	require.Equal(t, expectedName, user.Name())
	require.Equal(t, expectedHashedPassword, user.HashedPassword())
	require.ElementsMatch(t, expectedScopes, user.Scopes())
}

func TestNewUserProvider(t *testing.T) {
	provider := security.NewUserProvider()

	provider.Set("test", security.NewUser("test", []byte{}, []security.Scope{}))

	t.Run("user / exists", func(t *testing.T) {
		user := provider.Get("test")

		require.NotNil(t, user)
		require.Equal(t, "test", user.Name())
		require.Empty(t, user.HashedPassword())
		require.ElementsMatch(t, []security.Scope{}, user.Scopes())
	})

	t.Run("user / unknown", func(t *testing.T) {
		user := provider.Get("unknown")

		require.Nil(t, user)
	})
}
