package mock

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"time"
)

type Client struct{}

func (c *Client) Close() error {
	return nil
}

func (c *Client) Put(priority uint32, delay, ttr time.Duration, data []byte) (int, error) {
	return 1, nil
}

func (c *Client) Use(tube string) (string, error) {
	return tube, nil
}

func (c *Client) Reserve() (beanstalk.Job, error) {
	return beanstalk.Job{}, nil
}

func (c *Client) ReserveWithTimeout(timeout time.Duration) (beanstalk.Job, error) {
	return beanstalk.Job{}, nil
}

func (c *Client) ReserveJob(id int) (beanstalk.Job, error) {
	return beanstalk.Job{}, nil
}

func (c *Client) Delete(id int) error {
	switch id {
	case 999:
		return beanstalk.ErrNotFound

	default:
		return nil
	}
}

func (c *Client) Release(id int, priority uint32, delay time.Duration) error {
	switch id {
	case 999:
		return beanstalk.ErrNotFound

	default:
		return nil
	}
}

func (c *Client) Bury(id int, priority uint32) error {
	switch id {
	case 999:
		return beanstalk.ErrNotFound

	default:
		return nil
	}
}

func (c *Client) Touch(id int) error {
	return nil
}

func (c *Client) Watch(tube string) (int, error) {
	return 0, nil
}

func (c *Client) Ignore(tube string) (int, error) {
	return 0, nil
}

func (c *Client) Peek(id int) (beanstalk.Job, error) {
	switch id {
	case 999:
		return beanstalk.Job{}, beanstalk.ErrNotFound

	default:
		return beanstalk.Job{ID: id, Data: []byte("test")}, nil
	}
}

func (c *Client) PeekReady() (beanstalk.Job, error) {
	return beanstalk.Job{}, nil
}

func (c *Client) PeekDelayed() (beanstalk.Job, error) {
	return beanstalk.Job{}, nil
}

func (c *Client) PeekBuried() (beanstalk.Job, error) {
	return beanstalk.Job{}, nil
}

func (c *Client) Kick(bound int) (int, error) {
	return 0, nil
}

func (c *Client) KickJob(id int) error {
	switch id {
	case 999:
		return beanstalk.ErrNotFound

	default:
		return nil
	}
}

func (c *Client) StatsJob(id int) (beanstalk.StatsJob, error) {
	switch id {
	case 999:
		return beanstalk.StatsJob{}, beanstalk.ErrNotFound

	default:
		return beanstalk.StatsJob{}, nil
	}
}

func (c *Client) StatsTube(tube string) (beanstalk.StatsTube, error) {
	switch tube {
	case "not_found":
		return beanstalk.StatsTube{}, beanstalk.ErrNotFound

	default:
		return beanstalk.StatsTube{}, nil
	}
}

func (c *Client) Stats() (beanstalk.Stats, error) {
	return beanstalk.Stats{}, nil
}

func (c *Client) ListTubes() ([]string, error) {
	return []string{"default"}, nil
}

func (c *Client) ListTubeUsed() (string, error) {
	return "default", nil
}

func (c *Client) ListTubesWatched() ([]string, error) {
	return nil, nil
}

func (c *Client) PauseTube(tube string, delay time.Duration) error {
	return nil
}

func (c *Client) ExecuteCommand(command beanstalk.Command) (beanstalk.CommandResponse, error) {
	return nil, nil
}
