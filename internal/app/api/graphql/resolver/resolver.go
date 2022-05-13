package resolver

import "github.com/IvanLutokhin/go-beanstalk"

type ReleaseFunc func()

type Resolver struct {
	Pool beanstalk.Pool
}

func NewResolver(pool beanstalk.Pool) *Resolver {
	return &Resolver{Pool: pool}
}

func (r *Resolver) BeanstalkClient() (beanstalk.Client, ReleaseFunc) {
	client, err := r.Pool.Get()
	if err != nil {
		panic(err)
	}

	return client, func() { r.Pool.Put(client) }
}
