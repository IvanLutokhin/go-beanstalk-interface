package beanstalk

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
)

func NewPool(config *config.Config) beanstalk.Pool {
	pool, err := beanstalk.NewDefaultPool(config.Beanstalk.Address, config.Beanstalk.Pool.Capacity, false)
	if err != nil {
		panic(err)
	}

	return pool
}
