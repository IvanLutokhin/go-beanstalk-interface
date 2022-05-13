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

func TestQueryResolver_Server(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

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

func TestQueryResolver_ReadyJob(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

	c := client.New(h)

	q := `
query ReadyJob ($tube: String!) {
    readyJob(tube: $tube) {
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
		ReadyJob *model.Job
	}

	c.MustPost(q, &response, client.Var("tube", "default"))

	if id := response.ReadyJob.ID; id != 1 {
		t.Errorf("expected job id '1', but got '%v'", id)
	}
}

func TestQueryResolver_DelayedJob(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

	c := client.New(h)

	q := `
query DelayedJob ($tube: String!) {
    delayedJob(tube: $tube) {
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
		DelayedJob *model.Job
	}

	c.MustPost(q, &response, client.Var("tube", "default"))

	if id := response.DelayedJob.ID; id != 1 {
		t.Errorf("expected job id '1', but got '%v'", id)
	}
}

func TestQueryResolver_BuriedJob(t *testing.T) {
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: NewResolver(&mock.Pool{Client: &mock.Client{}})}))

	c := client.New(h)

	q := `
query BuriedJob ($tube: String!) {
    buriedJob(tube: $tube) {
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
		BuriedJob *model.Job
	}

	c.MustPost(q, &response, client.Var("tube", "default"))

	if id := response.BuriedJob.ID; id != 1 {
		t.Errorf("expected job id '1', but got '%v'", id)
	}
}
