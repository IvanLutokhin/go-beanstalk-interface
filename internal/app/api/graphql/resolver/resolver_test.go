package resolver_test

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/resolver"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/pkg/beanstalk/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResolver_AuthContext(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	r := resolver.NewResolver(pool)

	t.Run("user is not authenticated", func(t *testing.T) {
		err := r.AuthContext(context.Background(), []security.Scope{})

		require.NotNil(t, err)
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

		err := r.AuthContext(security.WithAuthenticatedUser(context.Background(), user), []security.Scope{security.ScopeWriteJobs})

		require.NotNil(t, err)
	})
}
