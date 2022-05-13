package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/dataloader"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
)

func (r *queryResolver) Server(ctx context.Context) (*model.Server, error) {
	client, release := r.BeanstalkClient()

	defer release()

	stats, err := client.Stats()
	if err != nil {
		return nil, err
	}

	return &model.Server{Stats: &stats}, nil
}

func (r *queryResolver) Tubes(ctx context.Context) (*model.TubeConnection, error) {
	client, release := r.BeanstalkClient()

	defer release()

	tubes, err := client.ListTubes()
	if err != nil {
		return nil, err
	}

	edges := make([]model.TubeEdge, len(tubes))

	for i, name := range tubes {
		node, err := dataloader.Tube(ctx, client, name)
		if err != nil {
			return nil, err
		}

		edges[i] = model.TubeEdge{Node: node}
	}

	return &model.TubeConnection{Edges: edges}, nil
}

func (r *queryResolver) Tube(ctx context.Context, name string) (*model.Tube, error) {
	client, release := r.BeanstalkClient()

	defer release()

	return dataloader.Tube(ctx, client, name)
}

func (r *queryResolver) Job(ctx context.Context, id int) (*model.Job, error) {
	client, release := r.BeanstalkClient()

	defer release()

	return dataloader.Job(ctx, client, id)
}

func (r *queryResolver) ReadyJob(ctx context.Context, tube string) (*model.Job, error) {
	client, release := r.BeanstalkClient()

	defer release()

	return dataloader.ReadyJob(ctx, client, tube)
}

func (r *queryResolver) DelayedJob(ctx context.Context, tube string) (*model.Job, error) {
	client, release := r.BeanstalkClient()

	defer release()

	return dataloader.DelayedJob(ctx, client, tube)
}

func (r *queryResolver) BuriedJob(ctx context.Context, tube string) (*model.Job, error) {
	client, release := r.BeanstalkClient()

	defer release()

	return dataloader.BuriedJob(ctx, client, tube)
}
