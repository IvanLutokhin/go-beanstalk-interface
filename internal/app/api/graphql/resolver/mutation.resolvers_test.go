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

func TestMutationResolver_BuryJob(t *testing.T) {
	mc := &mock.Client{}
	mc.On("Bury", 1, uint32(0)).Return(nil)
	mc.On("Bury", 999, uint32(100)).Return(beanstalk.ErrNotFound)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

	c := client.New(h)

	q := `
mutation BuryJob($id: Int!, $priority: Int!) {
   buryJob(input: {id: $id, priority: $priority}) {
       id
   }
}
`

	var response = struct {
		BuryJob *model.BuryJobPayload
	}{}

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
			client.Var("priority", 0),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.Equal(t, 1, response.BuryJob.ID)
	})

	t.Run("not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
			client.Var("priority", 100),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.NotNil(t, err)
		require.Equal(t, `[{"message":"beanstalk: not found","path":["buryJob"]}]`, err.Error())
	})
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

func TestMutationResolver_ReleaseJob(t *testing.T) {
	mc := &mock.Client{}
	mc.On("Release", 1, uint32(0), 5*time.Second).Return(nil)
	mc.On("Release", 999, uint32(100), 100*time.Second).Return(beanstalk.ErrNotFound)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

	c := client.New(h)

	q := `
mutation ReleaseJob($id: Int!, $priority: Int!, $delay: Int!) {
  releaseJob(input: {id: $id, priority: $priority, delay: $delay}) {
      id
  }
}
`
	var response = struct {
		ReleaseJob *model.ReleaseJobPayload
	}{}

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
			client.Var("priority", 0),
			client.Var("delay", 5*time.Second),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.Equal(t, 1, response.ReleaseJob.ID)
	})

	t.Run("not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
			client.Var("priority", 100),
			client.Var("delay", 100*time.Second),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs})),
		)

		require.NotNil(t, err)
		require.Equal(t, `[{"message":"beanstalk: not found","path":["releaseJob"]}]`, err.Error())
	})
}
