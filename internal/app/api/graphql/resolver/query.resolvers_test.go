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
)

func TestQueryResolver_Me(t *testing.T) {
	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(&mock.Pool{})}))

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

	require.Equal(t, "test", response.Me.User.Name)
	require.ElementsMatch(t, []model.Scope{model.ScopeReadServer}, response.Me.User.Scopes)
}

func TestQueryResolver_Server(t *testing.T) {
	expectedStats := beanstalk.Stats{
		CurrentJobsUrgent:     1,
		CurrentJobsReady:      1,
		CurrentJobsReserved:   1,
		CurrentJobsDelayed:    1,
		CurrentJobsBuried:     1,
		CmdPut:                1,
		CmdPeek:               1,
		CmdPeekReady:          1,
		CmdPeekDelayed:        1,
		CmdPeekBuried:         1,
		CmdReserve:            1,
		CmdUse:                1,
		CmdWatch:              1,
		CmdIgnore:             1,
		CmdDelete:             1,
		CmdRelease:            1,
		CmdBury:               1,
		CmdKick:               1,
		CmdTouch:              1,
		CmdStats:              1,
		CmdStatsJob:           1,
		CmdStatsTube:          1,
		CmdListTubes:          1,
		CmdListTubeUsed:       1,
		CmdListTubesWatched:   1,
		CmdPauseTube:          1,
		JobTimeouts:           10,
		TotalJobs:             25,
		MaxJobSize:            65535,
		CurrentTubes:          1,
		CurrentConnections:    1,
		CurrentProducers:      1,
		CurrentWorkers:        1,
		CurrentWaiting:        1,
		TotalConnections:      1,
		PID:                   1,
		Version:               "1.10",
		RUsageUTime:           0.148125,
		RUsageSTime:           0.014812,
		Uptime:                1864,
		BinlogOldestIndex:     1,
		BinlogCurrentIndex:    1,
		BinlogRecordsMigrated: 1,
		BinlogRecordsWritten:  1,
		BinlogMaxSize:         10485760,
		Draining:              false,
		ID:                    "f40521014b63360d",
		Hostname:              "671db3de0474",
	}

	mc := &mock.Client{}
	mc.On("Stats").Return(expectedStats, nil)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

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
           	cmdReserve,
           	cmdUse,
           	cmdWatch,
           	cmdIgnore,
           	cmdDelete,
           	cmdRelease,
           	cmdBury,
           	cmdKick,
			cmdTouch,
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

	require.Equal(t, expectedStats.CurrentJobsUrgent, response.Server.Stats.CurrentJobsUrgent)
	require.Equal(t, expectedStats.CurrentJobsReady, response.Server.Stats.CurrentJobsReady)
	require.Equal(t, expectedStats.CurrentJobsReserved, response.Server.Stats.CurrentJobsReserved)
	require.Equal(t, expectedStats.CurrentJobsDelayed, response.Server.Stats.CurrentJobsDelayed)
	require.Equal(t, expectedStats.CurrentJobsBuried, response.Server.Stats.CurrentJobsBuried)
	require.Equal(t, expectedStats.CmdPut, response.Server.Stats.CmdPut)
	require.Equal(t, expectedStats.CmdPeek, response.Server.Stats.CmdPeek)
	require.Equal(t, expectedStats.CmdPeekReady, response.Server.Stats.CmdPeekReady)
	require.Equal(t, expectedStats.CmdPeekDelayed, response.Server.Stats.CmdPeekDelayed)
	require.Equal(t, expectedStats.CmdPeekBuried, response.Server.Stats.CmdPeekBuried)
	require.Equal(t, expectedStats.CmdReserve, response.Server.Stats.CmdReserve)
	require.Equal(t, expectedStats.CmdUse, response.Server.Stats.CmdUse)
	require.Equal(t, expectedStats.CmdWatch, response.Server.Stats.CmdWatch)
	require.Equal(t, expectedStats.CmdIgnore, response.Server.Stats.CmdIgnore)
	require.Equal(t, expectedStats.CmdDelete, response.Server.Stats.CmdDelete)
	require.Equal(t, expectedStats.CmdRelease, response.Server.Stats.CmdRelease)
	require.Equal(t, expectedStats.CmdBury, response.Server.Stats.CmdBury)
	require.Equal(t, expectedStats.CmdKick, response.Server.Stats.CmdKick)
	require.Equal(t, expectedStats.CmdTouch, response.Server.Stats.CmdTouch)
	require.Equal(t, expectedStats.CmdStats, response.Server.Stats.CmdStats)
	require.Equal(t, expectedStats.CmdStatsJob, response.Server.Stats.CmdStatsJob)
	require.Equal(t, expectedStats.CmdStatsTube, response.Server.Stats.CmdStatsTube)
	require.Equal(t, expectedStats.CmdListTubes, response.Server.Stats.CmdListTubes)
	require.Equal(t, expectedStats.CmdListTubeUsed, response.Server.Stats.CmdListTubeUsed)
	require.Equal(t, expectedStats.CmdListTubesWatched, response.Server.Stats.CmdListTubesWatched)
	require.Equal(t, expectedStats.CmdPauseTube, response.Server.Stats.CmdPauseTube)
	require.Equal(t, expectedStats.JobTimeouts, response.Server.Stats.JobTimeouts)
	require.Equal(t, expectedStats.TotalJobs, response.Server.Stats.TotalJobs)
	require.Equal(t, expectedStats.MaxJobSize, response.Server.Stats.MaxJobSize)
	require.Equal(t, expectedStats.CurrentTubes, response.Server.Stats.CurrentTubes)
	require.Equal(t, expectedStats.CurrentConnections, response.Server.Stats.CurrentConnections)
	require.Equal(t, expectedStats.CurrentProducers, response.Server.Stats.CurrentProducers)
	require.Equal(t, expectedStats.CurrentWorkers, response.Server.Stats.CurrentWorkers)
	require.Equal(t, expectedStats.CurrentWaiting, response.Server.Stats.CurrentWaiting)
	require.Equal(t, expectedStats.TotalConnections, response.Server.Stats.TotalConnections)
	require.Equal(t, expectedStats.PID, response.Server.Stats.PID)
	require.Equal(t, expectedStats.Version, response.Server.Stats.Version)
	require.Equal(t, expectedStats.RUsageUTime, response.Server.Stats.RUsageUTime)
	require.Equal(t, expectedStats.RUsageSTime, response.Server.Stats.RUsageSTime)
	require.Equal(t, expectedStats.Uptime, response.Server.Stats.Uptime)
	require.Equal(t, expectedStats.BinlogOldestIndex, response.Server.Stats.BinlogOldestIndex)
	require.Equal(t, expectedStats.BinlogCurrentIndex, response.Server.Stats.BinlogCurrentIndex)
	require.Equal(t, expectedStats.BinlogRecordsMigrated, response.Server.Stats.BinlogRecordsMigrated)
	require.Equal(t, expectedStats.BinlogRecordsWritten, response.Server.Stats.BinlogRecordsWritten)
	require.Equal(t, expectedStats.BinlogMaxSize, response.Server.Stats.BinlogMaxSize)
	require.Equal(t, expectedStats.ID, response.Server.Stats.ID)
	require.Equal(t, expectedStats.Hostname, response.Server.Stats.Hostname)
}

func TestQueryResolver_Tubes(t *testing.T) {
	expectedTubeStats := beanstalk.StatsTube{
		Name:                "default",
		CurrentJobsUrgent:   1,
		CurrentJobsReady:    1,
		CurrentJobsReserved: 1,
		CurrentJobsDelayed:  1,
		CurrentJobsBuried:   1,
		TotalJobs:           5,
		CurrentUsing:        3,
		CurrentWaiting:      1,
		CurrentWatching:     2,
		Pause:               1,
		CmdDelete:           1,
		CmdPauseTube:        1,
		PauseTimeLeft:       10,
	}

	expectedReadyJobStats := beanstalk.StatsJob{
		ID:       1,
		Tube:     "default",
		State:    "ready",
		Priority: 999,
		Age:      12,
		Delay:    15,
		TTR:      1,
		TimeLeft: 10,
		File:     1,
		Reserves: 1,
		Timeouts: 1,
		Releases: 1,
		Buries:   1,
		Kicks:    1,
	}

	expectedDelayedJobStats := beanstalk.StatsJob{
		ID:       2,
		Tube:     "default",
		State:    "delayed",
		Priority: 0,
		Age:      12,
		Delay:    15,
		TTR:      1,
		TimeLeft: 10,
		File:     2,
		Reserves: 2,
		Timeouts: 2,
		Releases: 2,
		Buries:   2,
		Kicks:    2,
	}

	expectedBuriedJobStats := beanstalk.StatsJob{
		ID:       3,
		Tube:     "default",
		State:    "buried",
		Priority: 1,
		Age:      12,
		Delay:    15,
		TTR:      3,
		TimeLeft: 10,
		File:     3,
		Reserves: 4,
		Timeouts: 3,
		Releases: 3,
		Buries:   5,
		Kicks:    2,
	}

	mc := &mock.Client{}
	mc.On("ListTubes").Return([]string{"default"}, nil)
	mc.On("StatsTube", "default").Return(expectedTubeStats, nil)
	mc.On("Use", "default").Return("default", nil)
	mc.On("PeekReady").Return(beanstalk.Job{ID: 1, Data: []byte("test 1")}, nil)
	mc.On("StatsJob", 1).Return(expectedReadyJobStats, nil)
	mc.On("PeekDelayed").Return(beanstalk.Job{ID: 2, Data: []byte("test 2")}, nil)
	mc.On("StatsJob", 2).Return(expectedDelayedJobStats, nil)
	mc.On("PeekBuried").Return(beanstalk.Job{ID: 3, Data: []byte("test 3")}, nil)
	mc.On("StatsJob", 3).Return(expectedBuriedJobStats, nil)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

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

	require.Len(t, response.Tubes.Edges, 1)
	require.Equal(t, "default", response.Tubes.Edges[0].Node.Name)
	require.Equal(t, expectedTubeStats.CurrentJobsUrgent, response.Tubes.Edges[0].Node.Stats.CurrentJobsUrgent)
	require.Equal(t, expectedTubeStats.CurrentJobsReady, response.Tubes.Edges[0].Node.Stats.CurrentJobsReady)
	require.Equal(t, expectedTubeStats.CurrentJobsReserved, response.Tubes.Edges[0].Node.Stats.CurrentJobsReserved)
	require.Equal(t, expectedTubeStats.CurrentJobsDelayed, response.Tubes.Edges[0].Node.Stats.CurrentJobsDelayed)
	require.Equal(t, expectedTubeStats.CurrentJobsBuried, response.Tubes.Edges[0].Node.Stats.CurrentJobsBuried)
	require.Equal(t, expectedTubeStats.TotalJobs, response.Tubes.Edges[0].Node.Stats.TotalJobs)
	require.Equal(t, expectedTubeStats.CurrentUsing, response.Tubes.Edges[0].Node.Stats.CurrentUsing)
	require.Equal(t, expectedTubeStats.CurrentWaiting, response.Tubes.Edges[0].Node.Stats.CurrentWaiting)
	require.Equal(t, expectedTubeStats.CurrentWatching, response.Tubes.Edges[0].Node.Stats.CurrentWatching)
	require.Equal(t, expectedTubeStats.Pause, response.Tubes.Edges[0].Node.Stats.Pause)
	require.Equal(t, expectedTubeStats.CmdDelete, response.Tubes.Edges[0].Node.Stats.CmdDelete)
	require.Equal(t, expectedTubeStats.CmdPauseTube, response.Tubes.Edges[0].Node.Stats.CmdPauseTube)
	require.Equal(t, expectedTubeStats.PauseTimeLeft, response.Tubes.Edges[0].Node.Stats.PauseTimeLeft)

	require.Equal(t, 1, response.Tubes.Edges[0].Node.ReadyJob.ID)
	require.Equal(t, "test 1", response.Tubes.Edges[0].Node.ReadyJob.Data)
	require.Equal(t, expectedReadyJobStats.Tube, response.Tubes.Edges[0].Node.ReadyJob.Stats.Tube)
	require.Equal(t, expectedReadyJobStats.State, response.Tubes.Edges[0].Node.ReadyJob.Stats.State)
	require.Equal(t, expectedReadyJobStats.Priority, response.Tubes.Edges[0].Node.ReadyJob.Stats.Priority)
	require.Equal(t, expectedReadyJobStats.Age, response.Tubes.Edges[0].Node.ReadyJob.Stats.Age)
	require.Equal(t, expectedReadyJobStats.Delay, response.Tubes.Edges[0].Node.ReadyJob.Stats.Delay)
	require.Equal(t, expectedReadyJobStats.TTR, response.Tubes.Edges[0].Node.ReadyJob.Stats.TTR)
	require.Equal(t, expectedReadyJobStats.TimeLeft, response.Tubes.Edges[0].Node.ReadyJob.Stats.TimeLeft)
	require.Equal(t, expectedReadyJobStats.File, response.Tubes.Edges[0].Node.ReadyJob.Stats.File)
	require.Equal(t, expectedReadyJobStats.Reserves, response.Tubes.Edges[0].Node.ReadyJob.Stats.Reserves)
	require.Equal(t, expectedReadyJobStats.Timeouts, response.Tubes.Edges[0].Node.ReadyJob.Stats.Timeouts)
	require.Equal(t, expectedReadyJobStats.Releases, response.Tubes.Edges[0].Node.ReadyJob.Stats.Releases)
	require.Equal(t, expectedReadyJobStats.Buries, response.Tubes.Edges[0].Node.ReadyJob.Stats.Buries)
	require.Equal(t, expectedReadyJobStats.Kicks, response.Tubes.Edges[0].Node.ReadyJob.Stats.Kicks)

	require.Equal(t, 2, response.Tubes.Edges[0].Node.DelayedJob.ID)
	require.Equal(t, "test 2", response.Tubes.Edges[0].Node.DelayedJob.Data)
	require.Equal(t, expectedDelayedJobStats.Tube, response.Tubes.Edges[0].Node.DelayedJob.Stats.Tube)
	require.Equal(t, expectedDelayedJobStats.State, response.Tubes.Edges[0].Node.DelayedJob.Stats.State)
	require.Equal(t, expectedDelayedJobStats.Priority, response.Tubes.Edges[0].Node.DelayedJob.Stats.Priority)
	require.Equal(t, expectedDelayedJobStats.Age, response.Tubes.Edges[0].Node.DelayedJob.Stats.Age)
	require.Equal(t, expectedDelayedJobStats.Delay, response.Tubes.Edges[0].Node.DelayedJob.Stats.Delay)
	require.Equal(t, expectedDelayedJobStats.TTR, response.Tubes.Edges[0].Node.DelayedJob.Stats.TTR)
	require.Equal(t, expectedDelayedJobStats.TimeLeft, response.Tubes.Edges[0].Node.DelayedJob.Stats.TimeLeft)
	require.Equal(t, expectedDelayedJobStats.File, response.Tubes.Edges[0].Node.DelayedJob.Stats.File)
	require.Equal(t, expectedDelayedJobStats.Reserves, response.Tubes.Edges[0].Node.DelayedJob.Stats.Reserves)
	require.Equal(t, expectedDelayedJobStats.Timeouts, response.Tubes.Edges[0].Node.DelayedJob.Stats.Timeouts)
	require.Equal(t, expectedDelayedJobStats.Releases, response.Tubes.Edges[0].Node.DelayedJob.Stats.Releases)
	require.Equal(t, expectedDelayedJobStats.Buries, response.Tubes.Edges[0].Node.DelayedJob.Stats.Buries)
	require.Equal(t, expectedDelayedJobStats.Kicks, response.Tubes.Edges[0].Node.DelayedJob.Stats.Kicks)

	require.Equal(t, 3, response.Tubes.Edges[0].Node.BuriedJob.ID)
	require.Equal(t, "test 3", response.Tubes.Edges[0].Node.BuriedJob.Data)
	require.Equal(t, expectedBuriedJobStats.Tube, response.Tubes.Edges[0].Node.BuriedJob.Stats.Tube)
	require.Equal(t, expectedBuriedJobStats.State, response.Tubes.Edges[0].Node.BuriedJob.Stats.State)
	require.Equal(t, expectedBuriedJobStats.Priority, response.Tubes.Edges[0].Node.BuriedJob.Stats.Priority)
	require.Equal(t, expectedBuriedJobStats.Age, response.Tubes.Edges[0].Node.BuriedJob.Stats.Age)
	require.Equal(t, expectedBuriedJobStats.Delay, response.Tubes.Edges[0].Node.BuriedJob.Stats.Delay)
	require.Equal(t, expectedBuriedJobStats.TTR, response.Tubes.Edges[0].Node.BuriedJob.Stats.TTR)
	require.Equal(t, expectedBuriedJobStats.TimeLeft, response.Tubes.Edges[0].Node.BuriedJob.Stats.TimeLeft)
	require.Equal(t, expectedBuriedJobStats.File, response.Tubes.Edges[0].Node.BuriedJob.Stats.File)
	require.Equal(t, expectedBuriedJobStats.Reserves, response.Tubes.Edges[0].Node.BuriedJob.Stats.Reserves)
	require.Equal(t, expectedBuriedJobStats.Timeouts, response.Tubes.Edges[0].Node.BuriedJob.Stats.Timeouts)
	require.Equal(t, expectedBuriedJobStats.Releases, response.Tubes.Edges[0].Node.BuriedJob.Stats.Releases)
	require.Equal(t, expectedBuriedJobStats.Buries, response.Tubes.Edges[0].Node.BuriedJob.Stats.Buries)
	require.Equal(t, expectedBuriedJobStats.Kicks, response.Tubes.Edges[0].Node.BuriedJob.Stats.Kicks)
}

func TestQueryResolver_Tube(t *testing.T) {
	expectedTubeStats := beanstalk.StatsTube{
		Name:                "default",
		CurrentJobsUrgent:   1,
		CurrentJobsReady:    1,
		CurrentJobsReserved: 1,
		CurrentJobsDelayed:  1,
		CurrentJobsBuried:   1,
		TotalJobs:           5,
		CurrentUsing:        3,
		CurrentWaiting:      1,
		CurrentWatching:     2,
		Pause:               1,
		CmdDelete:           1,
		CmdPauseTube:        1,
		PauseTimeLeft:       10,
	}

	mc := &mock.Client{}
	mc.On("ListTubes").Return([]string{"default"}, nil)
	mc.On("StatsTube", "default").Return(expectedTubeStats, nil)
	mc.On("Use", "default").Return("default", nil)
	mc.On("StatsTube", "test").Return(beanstalk.StatsTube{}, beanstalk.ErrNotFound)
	mc.On("Use", "test").Return("test", nil)
	mc.On("PeekReady").Return(beanstalk.Job{}, beanstalk.ErrNotFound)
	mc.On("PeekDelayed").Return(beanstalk.Job{}, beanstalk.ErrNotFound)
	mc.On("PeekBuried").Return(beanstalk.Job{}, beanstalk.ErrNotFound)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

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

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("name", "default"),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs})),
		)

		require.Equal(t, "default", response.Tube.Name)
		require.Equal(t, expectedTubeStats.CurrentJobsUrgent, response.Tube.Stats.CurrentJobsUrgent)
		require.Equal(t, expectedTubeStats.CurrentJobsReady, response.Tube.Stats.CurrentJobsReady)
		require.Equal(t, expectedTubeStats.CurrentJobsReserved, response.Tube.Stats.CurrentJobsReserved)
		require.Equal(t, expectedTubeStats.CurrentJobsDelayed, response.Tube.Stats.CurrentJobsDelayed)
		require.Equal(t, expectedTubeStats.CurrentJobsBuried, response.Tube.Stats.CurrentJobsBuried)
		require.Equal(t, expectedTubeStats.TotalJobs, response.Tube.Stats.TotalJobs)
		require.Equal(t, expectedTubeStats.CurrentUsing, response.Tube.Stats.CurrentUsing)
		require.Equal(t, expectedTubeStats.CurrentWaiting, response.Tube.Stats.CurrentWaiting)
		require.Equal(t, expectedTubeStats.CurrentWatching, response.Tube.Stats.CurrentWatching)
		require.Equal(t, expectedTubeStats.Pause, response.Tube.Stats.Pause)
		require.Equal(t, expectedTubeStats.CmdDelete, response.Tube.Stats.CmdDelete)
		require.Equal(t, expectedTubeStats.CmdPauseTube, response.Tube.Stats.CmdPauseTube)
		require.Equal(t, expectedTubeStats.PauseTimeLeft, response.Tube.Stats.PauseTimeLeft)
	})

	t.Run("not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("name", "test"),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadTubes, security.ScopeReadJobs})),
		)

		require.NotNil(t, err)
		require.Equal(t, `[{"message":"beanstalk: not found","path":["tube","stats"]}]`, err.Error())
	})
}

func TestQueryResolver_Job(t *testing.T) {
	expectedStats := beanstalk.StatsJob{
		ID:       1,
		Tube:     "default",
		State:    "ready",
		Priority: 1,
		Age:      12,
		Delay:    15,
		TTR:      3,
		TimeLeft: 10,
		File:     3,
		Reserves: 4,
		Timeouts: 3,
		Releases: 3,
		Buries:   5,
		Kicks:    2,
	}

	mc := &mock.Client{}
	mc.On("Use", "default").Return("default", nil)
	mc.On("Peek", 1).Return(beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	mc.On("StatsJob", 1).Return(expectedStats, nil)
	mc.On("Peek", 999).Return(beanstalk.Job{}, beanstalk.ErrNotFound)

	h := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(mock.NewPool(mc))}))

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

	t.Run("success", func(t *testing.T) {
		c.MustPost(
			q,
			&response,
			client.Var("id", 1),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs})),
		)

		require.Equal(t, 1, response.Job.ID)
		require.Equal(t, "test", response.Job.Data)
		require.Equal(t, expectedStats.Tube, response.Job.Stats.Tube)
		require.Equal(t, expectedStats.State, response.Job.Stats.State)
		require.Equal(t, expectedStats.Priority, response.Job.Stats.Priority)
		require.Equal(t, expectedStats.Age, response.Job.Stats.Age)
		require.Equal(t, expectedStats.Delay, response.Job.Stats.Delay)
		require.Equal(t, expectedStats.TTR, response.Job.Stats.TTR)
		require.Equal(t, expectedStats.TimeLeft, response.Job.Stats.TimeLeft)
		require.Equal(t, expectedStats.File, response.Job.Stats.File)
		require.Equal(t, expectedStats.Reserves, response.Job.Stats.Reserves)
		require.Equal(t, expectedStats.Timeouts, response.Job.Stats.Timeouts)
		require.Equal(t, expectedStats.Releases, response.Job.Stats.Releases)
		require.Equal(t, expectedStats.Buries, response.Job.Stats.Buries)
		require.Equal(t, expectedStats.Kicks, response.Job.Stats.Kicks)
	})

	t.Run("not found", func(t *testing.T) {
		err := c.Post(
			q,
			&response,
			client.Var("id", 999),
			testutil.AuthenticatedUser(security.NewUser("test", []byte{}, []security.Scope{security.ScopeReadJobs})),
		)

		require.NotNil(t, err)
		require.Equal(t, `[{"message":"beanstalk: not found","path":["job"]}]`, err.Error())
	})
}
