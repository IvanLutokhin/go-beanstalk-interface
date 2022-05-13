package dataloader

import (
	"context"
	"errors"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
)

func Tube(ctx context.Context, client beanstalk.Client, name string) (*model.Tube, error) {
	stats, err := client.StatsTube(name)
	if err != nil {
		return nil, err
	}

	readyJob, err := ReadyJob(ctx, client, name)
	if err != nil && !errors.Is(err, beanstalk.ErrNotFound) {
		return nil, err
	}

	delayedJob, err := DelayedJob(ctx, client, name)
	if err != nil && !errors.Is(err, beanstalk.ErrNotFound) {
		return nil, err
	}

	buriedJob, err := BuriedJob(ctx, client, name)
	if err != nil && !errors.Is(err, beanstalk.ErrNotFound) {
		return nil, err
	}

	tube := &model.Tube{
		Name:       name,
		Stats:      &stats,
		ReadyJob:   readyJob,
		DelayedJob: delayedJob,
		BuriedJob:  buriedJob,
	}

	return tube, nil
}
