package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
)

func (r *userResolver) Scopes(ctx context.Context, obj *security.User) ([]model.Scope, error) {
	scopes := make([]model.Scope, len(obj.Scopes()))
	for i, scope := range obj.Scopes() {
		s, err := model.MapScope(scope)
		if err != nil {
			return nil, err
		}

		scopes[i] = s
	}

	return scopes, nil
}

// User returns executor.UserResolver implementation.
func (r *Resolver) User() executor.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
