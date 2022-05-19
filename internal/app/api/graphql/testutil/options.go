package testutil

import (
	"github.com/99designs/gqlgen/client"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
)

func AuthenticatedUser(user *security.User) client.Option {
	return func(bd *client.Request) {
		ctx := security.WithAuthenticatedUser(bd.HTTP.Context(), user)

		bd.HTTP = bd.HTTP.WithContext(ctx)
	}
}
