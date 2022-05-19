package resolver

import (
	"context"
	"errors"
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
)

type ReleaseFunc func()

type Resolver struct {
	Pool beanstalk.Pool
}

func NewResolver(pool beanstalk.Pool) *Resolver {
	return &Resolver{Pool: pool}
}

func (r *Resolver) BeanstalkClient() (beanstalk.Client, ReleaseFunc) {
	client, err := r.Pool.Get()
	if err != nil {
		panic(err)
	}

	return client, func() { r.Pool.Put(client) }
}

func (r *Resolver) AuthContext(ctx context.Context, expectedScopes []security.Scope) error {
	user := security.AuthenticatedUser(ctx)
	if user == nil {
		return errors.New("user is not authenticated")
	}

	if scopes := security.VerifyScopes(user.Scopes(), expectedScopes); len(scopes) > 0 {
		return fmt.Errorf("required scopes %v", scopes)
	}

	return nil
}
