package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
)

func (r *mutationResolver) CreateJob(ctx context.Context, input *model.CreateJobInput) (*model.CreateJobPayload, error) {
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

func (r *mutationResolver) BuryJob(ctx context.Context, input *model.BuryJobInput) (*model.BuryJobPayload, error) {
	client, release := r.BeanstalkClient()

	defer release()

	err := client.Bury(input.ID, uint32(input.Priority))
	if err != nil {
		return nil, err
	}

	return &model.BuryJobPayload{ID: input.ID}, nil
}

func (r *mutationResolver) DeleteJob(ctx context.Context, input *model.DeleteJobInput) (*model.DeleteJobPayload, error) {
	client, release := r.BeanstalkClient()

	defer release()

	err := client.Delete(input.ID)
	if err != nil {
		return nil, err
	}

	return &model.DeleteJobPayload{ID: input.ID}, nil
}

func (r *mutationResolver) KickJob(ctx context.Context, input *model.KickJobInput) (*model.KickJobPayload, error) {
	client, release := r.BeanstalkClient()

	defer release()

	err := client.KickJob(input.ID)
	if err != nil {
		return nil, err
	}

	return &model.KickJobPayload{ID: input.ID}, nil
}

func (r *mutationResolver) ReleaseJob(ctx context.Context, input *model.ReleaseJobInput) (*model.ReleaseJobPayload, error) {
	client, release := r.BeanstalkClient()

	defer release()

	err := client.Release(input.ID, uint32(input.Priority), time.Duration(input.Delay))
	if err != nil {
		return nil, err
	}

	return &model.ReleaseJobPayload{ID: input.ID}, nil
}
