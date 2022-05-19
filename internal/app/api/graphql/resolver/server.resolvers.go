package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"

	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
)

func (r *serverResolver) Stats(ctx context.Context, obj *model.Server) (*beanstalk.Stats, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadServer}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	stats, err := client.Stats()
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// Server returns executor.ServerResolver implementation.
func (r *Resolver) Server() executor.ServerResolver { return &serverResolver{r} }

type serverResolver struct{ *Resolver }
