package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/model"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/pkg/beanstalk/mock"
	"strings"
	"testing"
)

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

	c.MustPost(q, &response)
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

	c.MustPost(q, &response)

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

	c.MustPost(q, &response, client.Var("name", "default"))

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

	c.MustPost(q, &response, client.Var("id", 1))

	if id := response.Job.ID; id != 1 {
		t.Errorf("expected job id '1', but got '%v'", id)
	}
}
