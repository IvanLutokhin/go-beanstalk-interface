package beanstalk

import (
	"context"
	"github.com/IvanLutokhin/go-beanstalk"
	"go.uber.org/fx"
)

func RegisterHooks(lifecycle fx.Lifecycle, pool beanstalk.Pool) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return pool.Open(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return pool.Close(ctx)
		},
	})
}
