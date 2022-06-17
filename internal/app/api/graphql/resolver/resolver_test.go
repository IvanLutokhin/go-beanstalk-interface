package resolver_test

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/resolver"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResolver_AuthContext(t *testing.T) {
	r := resolver.NewResolver(&mock.Pool{})

	t.Run("user is not authenticated", func(t *testing.T) {
		require.Error(t, r.AuthContext(context.Background(), []security.Scope{}))
	})

	t.Run("missing scopes", func(t *testing.T) {
		user := security.NewUser(
			"test",
			[]byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC"),
			[]security.Scope{
				security.ScopeReadServer,
				security.ScopeReadTubes,
				security.ScopeReadJobs,
			},
		)

		require.Error(t, r.AuthContext(security.WithAuthenticatedUser(context.Background(), user), []security.Scope{security.ScopeWriteJobs}))
	})
}
