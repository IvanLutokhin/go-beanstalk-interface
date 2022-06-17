package beanstalk

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
)

func NewPool(config *config.Config, logger *LoggerAdapter) beanstalk.Pool {
	options := &beanstalk.PoolOptions{
		Dialer: func() (*beanstalk.DefaultClient, error) {
			return beanstalk.Dial(config.Beanstalk.Address)
		},
		Logger:      logger,
		Capacity:    config.Beanstalk.Pool.Capacity,
		MaxAge:      config.Beanstalk.Pool.MaxAge,
		IdleTimeout: config.Beanstalk.Pool.IdleTimeout,
	}

	return beanstalk.NewDefaultPool(options)
}
