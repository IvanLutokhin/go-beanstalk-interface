package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
)

func (r *jobResolver) Stats(ctx context.Context, obj *model.Job) (*beanstalk.StatsJob, error) {
	client, release := r.BeanstalkClient()

	defer release()

	stats, err := client.StatsJob(obj.ID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// Job returns executor.JobResolver implementation.
func (r *Resolver) Job() executor.JobResolver { return &jobResolver{r} }

type jobResolver struct{ *Resolver }
