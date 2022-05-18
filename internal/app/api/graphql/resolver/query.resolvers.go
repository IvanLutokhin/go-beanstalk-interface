package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
)

func (r *queryResolver) Server(ctx context.Context) (*model.Server, error) {
	return &model.Server{}, nil
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
		edges[i] = model.TubeEdge{
			Node: &model.Tube{
				Name: name,
			},
		}
	}

	return &model.TubeConnection{Edges: edges}, nil
}

func (r *queryResolver) Tube(ctx context.Context, name string) (*model.Tube, error) {
	return &model.Tube{Name: name}, nil
}

func (r *queryResolver) Job(ctx context.Context, id int) (*model.Job, error) {
	client, release := r.BeanstalkClient()

	defer release()

	peeked, err := client.Peek(id)
	if err != nil {
		return nil, err
	}

	job := model.Job{
		ID:   peeked.ID,
		Data: string(peeked.Data),
	}

	return &job, nil
}
