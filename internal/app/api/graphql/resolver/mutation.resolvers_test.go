package resolver_test

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/resolver"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/testutil"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMutationResolver_CreateJob(t *testing.T) {
	mc := &mock.Client{}
	mc.On("Use", "default").Return("default", nil)
	mc.On("Put", uint32(0), time.Duration(0), time.Duration(0), []byte("test")).Return(1, nil)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

	c := client.New(h)

	q := `
mutation CreateJob($tube: String!, $priority: Int!, $delay: Int!, $ttr: Int!, $data: String!) {
	createJob(input: {tube: $tube, priority: $priority, delay: $delay, ttr: $ttr, data: $data}) {
		tube,
		id
	}
}
`
	var response = struct {
		CreateJob *model.CreateJobPayload
	}{}

	c.MustPost(
		q,
		&response,
		client.Var("tube", "default"),
		client.Var("priority", 0),
		client.Var("delay", 0),
		client.Var("ttr", 0),
		client.Var("data", "test"),
		testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
	)

	require.Equal(t, "default", response.CreateJob.Tube)
	require.Equal(t, 1, response.CreateJob.ID)
}

func TestMutationResolver_DeleteJob(t *testing.T) {
	mc := &mock.Client{}
	mc.On("Delete", 1).Return(nil)
	mc.On("Delete", 999).Return(beanstalk.ErrNotFound)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

	c := client.New(h)

	q := `
mutation DeleteJob($id: Int!) {
  deleteJob(input: {id: $id}) {
      id
  }
}
`
	var response = struct {
		DeleteJob *model.DeleteJobPayload
	}{}

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.Equal(t, 1, response.DeleteJob.ID)
	})

	t.Run("not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.NotNil(t, err)
		require.Equal(t, `[{"message":"beanstalk: not found","path":["deleteJob"]}]`, err.Error())
	})
}

func TestMutationResolver_DeleteJobs(t *testing.T) {
	mc := &mock.Client{}
	mc.On("Use", "default").Return("default", nil)
	mc.On("PeekBuried").Times(5).Return(&beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	mc.On("PeekBuried").Return(nil, beanstalk.ErrNotFound)
	mc.On("Delete", 1).Return(nil)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

	c := client.New(h)

	q := `
mutation DeleteJobs($tube: String!, $count: Int) {
  deleteJobs(input: {tube: $tube, count: $count}) {
      count
  }
}
`
	var response = struct {
		DeleteJobs *model.DeleteJobsPayload
	}{}

	t.Run("with limit", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("tube", "default"),
			client.Var("count", 3),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.Equal(t, 3, response.DeleteJobs.Count)
	})

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("tube", "default"),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.Equal(t, 2, response.DeleteJobs.Count)
	})
}

func TestMutationResolver_KickJob(t *testing.T) {
	mc := &mock.Client{}
	mc.On("KickJob", 1).Return(nil)
	mc.On("KickJob", 999).Return(beanstalk.ErrNotFound)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

	c := client.New(h)

	q := `
mutation KickJob($id: Int!) {
  kickJob(input: {id: $id}) {
      id
  }
}
`
	var response = struct {
		KickJob *model.KickJobPayload
	}{}

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.Equal(t, 1, response.KickJob.ID)
	})

	t.Run("not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.NotNil(t, err)
		require.Equal(t, `[{"message":"beanstalk: not found","path":["kickJob"]}]`, err.Error())
	})
}

func TestMutationResolver_KickJobs(t *testing.T) {
	mc := &mock.Client{}
	mc.On("Use", "default").Return("default", nil)
	mc.On("PeekBuried").Times(7).Return(&beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	mc.On("PeekBuried").Return(nil, beanstalk.ErrNotFound)
	mc.On("KickJob", 1).Return(nil)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

	c := client.New(h)

	q := `
mutation KickJobs($tube: String!, $count: Int) {
  kickJobs(input: {tube: $tube, count: $count}) {
      count
  }
}
`
	var response = struct {
		KickJobs *model.KickJobsPayload
	}{}

	t.Run("with limit", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("tube", "default"),
			client.Var("count", 5),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.Equal(t, 5, response.KickJobs.Count)
	})

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("tube", "default"),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.Equal(t, 2, response.KickJobs.Count)
	})
}

func TestMutationResolver_PauseTube(t *testing.T) {
	mc := &mock.Client{}
	mc.On("PauseTube", "default", 30*time.Minute).Return(nil)
	mc.On("PauseTube", "test", 5*time.Second).Return(beanstalk.ErrNotFound)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

	c := client.New(h)

	q := `
mutation PauseTube($tube: String!, $delay: Int!) {
  pauseTube(input: {tube: $tube, delay: $delay}) {
      tube
  }
}
`
	var response = struct {
		PauseTube *model.PauseTubePayload
	}{}

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("tube", "default"),
			client.Var("delay", 30*time.Minute),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeWriteTubes})),
		)

		require.Equal(t, "default", response.PauseTube.Tube)
	})

	t.Run("not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("tube", "test"),
			client.Var("delay", 5*time.Second),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeWriteTubes})),
		)

		require.NotNil(t, err)
		require.Equal(t, `[{"message":"beanstalk: not found","path":["pauseTube"]}]`, err.Error())
	})
}
