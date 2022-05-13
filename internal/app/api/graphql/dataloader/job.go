package dataloader

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
)

func Job(ctx context.Context, client beanstalk.Client, id int) (*model.Job, error) {
	peeked, err := client.Peek(id)
	if err != nil {
		return nil, err
	}

	stats, err := client.StatsJob(peeked.ID)
	if err != nil {
		return nil, err
	}

	job := &model.Job{
		ID:    peeked.ID,
		Data:  string(peeked.Data),
		Stats: &stats,
	}

	return job, nil
}

func ReadyJob(ctx context.Context, client beanstalk.Client, tube string) (*model.Job, error) {
	_, err := client.Use(tube)
	if err != nil {
		return nil, err
	}

	peeked, err := client.PeekReady()
	if err != nil {
		return nil, err
	}

	stats, err := client.StatsJob(peeked.ID)
	if err != nil {
		return nil, err
	}

	job := &model.Job{
		ID:    peeked.ID,
		Data:  string(peeked.Data),
		Stats: &stats,
	}

	return job, nil
}

func DelayedJob(ctx context.Context, client beanstalk.Client, tube string) (*model.Job, error) {
	_, err := client.Use(tube)
	if err != nil {
		return nil, err
	}

	peeked, err := client.PeekDelayed()
	if err != nil {
		return nil, err
	}

	stats, err := client.StatsJob(peeked.ID)
	if err != nil {
		return nil, err
	}

	job := &model.Job{
		ID:    peeked.ID,
		Data:  string(peeked.Data),
		Stats: &stats,
	}

	return job, nil
}

func BuriedJob(ctx context.Context, client beanstalk.Client, tube string) (*model.Job, error) {
	_, err := client.Use(tube)
	if err != nil {
		return nil, err
	}

	peeked, err := client.PeekBuried()
	if err != nil {
		return nil, err
	}

	stats, err := client.StatsJob(peeked.ID)
	if err != nil {
		return nil, err
	}

	job := &model.Job{
		ID:    peeked.ID,
		Data:  string(peeked.Data),
		Stats: &stats,
	}

	return job, nil
}
