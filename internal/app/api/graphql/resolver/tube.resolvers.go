package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
)

func (r *tubeResolver) Stats(ctx context.Context, obj *model.Tube) (*beanstalk.StatsTube, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadTubes}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	stats, err := client.StatsTube(obj.Name)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *tubeResolver) ReadyJob(ctx context.Context, obj *model.Tube) (*model.Job, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadJobs}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	_, err := client.Use(obj.Name)
	if err != nil {
		return nil, err
	}

	peeked, err := client.PeekReady()
	if err != nil {
		if errors.Is(err, beanstalk.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	job := model.Job{
		ID:   peeked.ID,
		Data: string(peeked.Data),
	}

	return &job, nil
}

func (r *tubeResolver) DelayedJob(ctx context.Context, obj *model.Tube) (*model.Job, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadJobs}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	_, err := client.Use(obj.Name)
	if err != nil {
		return nil, err
	}

	peeked, err := client.PeekDelayed()
	if err != nil {
		if errors.Is(err, beanstalk.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	job := model.Job{
		ID:   peeked.ID,
		Data: string(peeked.Data),
	}

	return &job, nil
}

func (r *tubeResolver) BuriedJob(ctx context.Context, obj *model.Tube) (*model.Job, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadJobs}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	_, err := client.Use(obj.Name)
	if err != nil {
		return nil, err
	}

	peeked, err := client.PeekBuried()
	if err != nil {
		if errors.Is(err, beanstalk.ErrNotFound) {
			return nil, nil
		}

		return nil, err
	}

	job := model.Job{
		ID:   peeked.ID,
		Data: string(peeked.Data),
	}

	return &job, nil
}

// Tube returns executor.TubeResolver implementation.
func (r *Resolver) Tube() executor.TubeResolver { return &tubeResolver{r} }

type tubeResolver struct{ *Resolver }
