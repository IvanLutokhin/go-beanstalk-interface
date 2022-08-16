package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
)

func (r *mutationResolver) CreateJob(ctx context.Context, input *model.CreateJobInput) (*model.CreateJobPayload, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	tube, err := client.Use(input.Tube)
	if err != nil {
		return nil, err
	}

	id, err := client.Put(uint32(input.Priority), time.Duration(input.Delay), time.Duration(input.Ttr), []byte(input.Data))
	if err != nil {
		return nil, err
	}

	return &model.CreateJobPayload{Tube: tube, ID: id}, nil
}

func (r *mutationResolver) DeleteJob(ctx context.Context, input *model.DeleteJobInput) (*model.DeleteJobPayload, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	err := client.Delete(input.ID)
	if err != nil {
		return nil, err
	}

	return &model.DeleteJobPayload{ID: input.ID}, nil
}

func (r *mutationResolver) DeleteJobs(ctx context.Context, input *model.DeleteJobsInput) (*model.DeleteJobsPayload, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs, security.ScopeWriteJobs}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	if _, err := client.Use(input.Tube); err != nil {
		return nil, err
	}

	count := input.Count

	payload := &model.DeleteJobsPayload{Count: 0}

	for {
		peeked, err := client.PeekBuried()
		if err != nil && !errors.Is(err, beanstalk.ErrNotFound) {
			return payload, err
		}

		if peeked == nil {
			break
		}

		if err := client.Delete(peeked.ID); err != nil {
			return payload, err
		}

		payload.Count++

		if count != nil && *count <= payload.Count {
			break
		}
	}

	return payload, nil
}

func (r *mutationResolver) KickJob(ctx context.Context, input *model.KickJobInput) (*model.KickJobPayload, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	err := client.KickJob(input.ID)
	if err != nil {
		return nil, err
	}

	return &model.KickJobPayload{ID: input.ID}, nil
}

func (r *mutationResolver) KickJobs(ctx context.Context, input *model.KickJobsInput) (*model.KickJobsPayload, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs, security.ScopeWriteJobs}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	if _, err := client.Use(input.Tube); err != nil {
		return nil, err
	}

	count := input.Count

	payload := &model.KickJobsPayload{Count: 0}

	for {
		peeked, err := client.PeekBuried()
		if err != nil && !errors.Is(err, beanstalk.ErrNotFound) {
			return payload, err
		}

		if peeked == nil {
			break
		}

		if err := client.KickJob(peeked.ID); err != nil {
			return payload, err
		}

		payload.Count++

		if count != nil && *count <= payload.Count {
			break
		}
	}

	return payload, nil
}

func (r *mutationResolver) PauseTube(ctx context.Context, input *model.PauseTubeInput) (*model.PauseTubePayload, error) {
	if err := r.AuthContext(ctx, []security.Scope{security.ScopeReadTubes, security.ScopeWriteTubes}); err != nil {
		return nil, err
	}

	client, release := r.BeanstalkClient()

	defer release()

	err := client.PauseTube(input.Tube, time.Duration(input.Delay))
	if err != nil {
		return nil, err
	}

	return &model.PauseTubePayload{Tube: input.Tube}, nil
}
