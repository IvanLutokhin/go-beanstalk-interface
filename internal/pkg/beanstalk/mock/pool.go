package mock

import "github.com/IvanLutokhin/go-beanstalk"

type Pool struct {
	Client beanstalk.Client
}

func (p *Pool) Open() error {
	return nil
}

func (p *Pool) Close() error {
	return nil
}

func (p *Pool) Get() (beanstalk.Client, error) {
	return p.Client, nil
}

func (p *Pool) Put(client beanstalk.Client) error {
	p.Client = client

	return nil
}

func (p *Pool) Len() int {
	return 1
}
