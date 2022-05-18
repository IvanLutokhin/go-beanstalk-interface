package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
)

// Mutation returns executor.MutationResolver implementation.
func (r *Resolver) Mutation() executor.MutationResolver { return &mutationResolver{r} }

// Query returns executor.QueryResolver implementation.
func (r *Resolver) Query() executor.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
