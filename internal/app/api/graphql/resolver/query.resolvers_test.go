package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/testutil"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/pkg/beanstalk/mock"
	"reflect"
	"strings"
	"testing"
)

func TestQueryResolver_Me(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 3, true)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: NewResolver(pool)}))

	c := client.New(h)

	q := `
query Me() {
	me {
		user {
			name,
			scopes
		}
	}
}
`

	var response = struct {
		Me struct {
			User struct {
				Name   string
				Scopes []model.Scope
			}
		}
	}{}

	c.MustPost(
		q,
		&response,
		testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadServer})),
	)

	if name := response.Me.User.Name; !strings.EqualFold("test", name) {
		t.Errorf("expected user name 'test', but got '%v'", name)
	}

	if scopes := response.Me.User.Scopes; !reflect.DeepEqual([]model.Scope{model.ScopeReadServer}, scopes) {
		t.Errorf("expected user scopes '%v', but got '%v'", []model.Scope{model.ScopeReadServer}, scopes)
	}
}

func TestQueryResolver_Server(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 3, true)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: NewResolver(pool)}))

	c := client.New(h)

	q := `
query Server() {
	server {
		stats {
			currentJobsUrgent,
            currentJobsReady,
            currentJobsReserved,
            currentJobsDelayed,
            currentJobsBuried,
            cmdPut,
            cmdPeek,
            cmdPeekReady,
            cmdPeekDelayed,
            cmdPeekBuried,
            cmdPeekReserve,
            cmdPeekUse,
            cmdWatch,
            cmdIgnore,
            cmdDelete,
            cmdRelease,
            cmdBury,
            cmdKick,
            cmdStats,
            cmdStatsJob,
            cmdStatsTube,
            cmdListTubes,
            cmdListTubeUsed,
            cmdListTubesWatched,
            cmdPauseTube,
            jobTimeouts,
            totalJobs,
            maxJobSize,
            currentTubes,
            currentConnections,
            currentProducers,
            currentWorkers,
            currentWaiting,
            totalConnections,
            pid,
            version,
            rUsageUTime,
            rUsageSTime,
            uptime,
            binlogOldestIndex,
            binlogCurrentIndex,
            binlogMaxSize,
            binlogRecordsWritten,
            binlogRecordsMigrated,
            draining,
            id,
            hostname,
            os,
            platform
		}
	}
}
`

	var response = struct {
		Server *model.Server
	}{}

	c.MustPost(
		q,
		&response,
		testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadServer})),
	)
}

func TestQueryResolver_Tubes(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 3, true)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: NewResolver(pool)}))

	c := client.New(h)

	q := `
query TubeList {
   tubes {
       edges {
           node {
               name,
               stats {
                   currentJobsUrgent,
                   currentJobsReady,
                   currentJobsReserved,
                   currentJobsDelayed,
                   currentJobsBuried,
                   totalJobs,
                   currentUsing,
                   currentWaiting,
                   currentWatching,
                   pause,
                   cmdDelete,
                   cmdPauseTube,
                   pauseTimeLeft
               },
               readyJob {
                   ...job
               },
               delayedJob {
                   ...job
               },
               buriedJob {
                   ...job
               }
           }
       }
   }
}

fragment job on Job {
   id,
   data,
   stats {
       tube,
       state,
       priority,
       age,
       delay,
       ttr,
       timeLeft,
       file,
       reserves,
       timeouts,
       releases,
       buries,
       kicks
   }
}
`

	var response struct {
		Tubes struct {
			Edges []model.TubeEdge
		}
	}

	c.MustPost(
		q,
		&response,
		testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs})),
	)

	if count := len(response.Tubes.Edges); count != 1 {
		t.Errorf("expected tubes '1', but got '%v'", count)
	}
}

func TestQueryResolver_Tube(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 3, true)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: NewResolver(pool)}))

	c := client.New(h)

	q := `
query Tube ($name: String!) {
    tube(name: $name) {
        name,
        stats {
            currentJobsUrgent,
            currentJobsReady,
            currentJobsReserved,
            currentJobsDelayed,
            currentJobsBuried,
            totalJobs,
            currentUsing,
            currentWaiting,
            currentWatching,
            pause,
            cmdDelete,
            cmdPauseTube,
            pauseTimeLeft
        },
        readyJob {
            ...job
        },
        delayedJob {
            ...job
        },
        buriedJob {
            ...job
        }
    }
}

fragment job on Job {
    id,
    data,
    stats {
        tube,
        state,
        priority,
        age,
        delay,
        ttr,
        timeLeft,
        file,
        reserves,
        timeouts,
        releases,
        buries,
        kicks
    }
}
`

	var response struct {
		Tube *model.Tube
	}

	c.MustPost(
		q,
		&response,
		client.Var("name", "default"),
		testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs})),
	)

	if name := response.Tube.Name; !strings.EqualFold("default", name) {
		t.Errorf("expected tube name 'default', but got '%v'", name)
	}
}

func TestQueryResolver_Job(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 3, true)
	if err != nil {
		t.Fatal(err)
	}

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: NewResolver(pool)}))

	c := client.New(h)

	q := `
query Job ($id: Int!) {
    job(id: $id) {
		id,
		data,
		stats {
			tube,
			state,
			priority,
			age,
			delay,
			ttr,
			timeLeft,
			file,
			reserves,
			timeouts,
			releases,
			buries,
			kicks
		}
	}
}
`

	var response struct {
		Job *model.Job
	}

	c.MustPost(
		q,
		&response,
		client.Var("id", 1),
		testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs})),
	)

	if id := response.Job.ID; id != 1 {
		t.Errorf("expected job id '1', but got '%v'", id)
	}
}
