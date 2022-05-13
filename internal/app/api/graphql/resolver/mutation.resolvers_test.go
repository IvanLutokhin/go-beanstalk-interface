package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/pkg/beanstalk/mock"
	"strings"
	"testing"
)

func TestMutationResolver_CreateJob(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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
	)

	if tube := response.CreateJob.Tube; !strings.EqualFold("default", tube) {
		t.Errorf("expected tube name 'default', but got '%v'", tube)
	}

	if id := response.CreateJob.ID; id != 1 {
		t.Errorf("expected job id '1', but got '%v'", id)
	}
}

func TestMutationResolver_BuryJob(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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

	t.Run("bury job", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
			client.Var("priority", 0),
		)

		if id := response.BuryJob.ID; id != 1 {
			t.Errorf("expected job id '1', but got '%v'", id)
		}
	})

	t.Run("bury job / not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
			client.Var("priority", 0),
		)

		if !strings.EqualFold(`[{"message":"beanstalk: not found","path":["buryJob"]}]`, err.Error()) {
			t.Errorf("expected not found error, but got '%v'", err)
		}
	})
}

func TestMutationResolver_DeleteJob(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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

	t.Run("delete job", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
		)

		if id := response.DeleteJob.ID; id != 1 {
			t.Errorf("expected job id '1', but got '%v'", id)
		}
	})

	t.Run("delete job / not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
		)

		if !strings.EqualFold(`[{"message":"beanstalk: not found","path":["deleteJob"]}]`, err.Error()) {
			t.Errorf("expected not found error, but got '%v'", err)
		}
	})
}

func TestMutationResolver_KickJob(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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

	t.Run("kick job", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
		)

		if id := response.KickJob.ID; id != 1 {
			t.Errorf("expected job id '1', but got '%v'", id)
		}
	})

	t.Run("kick job / not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
		)

		if !strings.EqualFold(`[{"message":"beanstalk: not found","path":["kickJob"]}]`, err.Error()) {
			t.Errorf("expected not found error, but got '%v'", err)
		}
	})
}

func TestMutationResolver_ReleaseJob(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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

	t.Run("release job", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
			client.Var("priority", 0),
			client.Var("delay", 0),
		)

		if id := response.ReleaseJob.ID; id != 1 {
			t.Errorf("expected job id '1', but got '%v'", id)
		}
	})

	t.Run("release job / not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
			client.Var("priority", 0),
			client.Var("delay", 0),
		)

		if !strings.EqualFold(`[{"message":"beanstalk: not found","path":["releaseJob"]}]`, err.Error()) {
			t.Errorf("expected not found error, but got '%v'", err)
		}
	})
}
