package handler_test

import (
	"bytes"
	"flag"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/cli/handler"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"strings"
	"testing"
	"time"
)

func TestHandler_Put(t *testing.T) {
	t.Run("empty argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("tube", "default", "")
		flagSet.Int("priority", 0, "")
		flagSet.Duration("delay", 0, "")
		flagSet.Duration("ttr", 0, "")
		flagSet.String("format", "json", "")

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.Error(t, h.Put(ctx))
	})

	t.Run("success", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Use", "default").Return("default", nil)
		client.On("Put", uint32(0), time.Duration(0), time.Duration(0), []byte("test")).Return(1, nil)
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("tube", "default", "")
		flagSet.Int("priority", 0, "")
		flagSet.Duration("delay", 0, "")
		flagSet.Duration("ttr", 0, "")
		flagSet.String("format", "json", "")
		flagSet.Parse([]string{"test"})

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.Put(ctx))

		var response handler.PutCommandResponse
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, "default", response.Tube)
		require.Equal(t, 1, response.ID)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("empty argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.Delete(ctx))
	})

	t.Run("invalid argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.Parse([]string{"test"})

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.Delete(ctx))
	})

	t.Run("success", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Delete", 1).Return(nil)
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.Parse([]string{"1"})

		ctx := cli.NewContext(nil, flagSet, nil)

		require.NoError(t, h.Delete(ctx))
	})
}

func TestHandler_DeleteJobs(t *testing.T) {
	client := &mock.Client{}
	client.On("Use", "default").Return("default", nil)
	client.On("PeekBuried").Times(5).Return(&beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	client.On("PeekBuried").Return(nil, beanstalk.ErrNotFound)
	client.On("Delete", 1).Return(nil)
	client.On("Close", mock.Anything).Return(nil)

	h := &handler.Handler{
		Client: client,
	}

	t.Run("with limit", func(t *testing.T) {
		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("tube", "default", "")
		flagSet.Int("count", 3, "")
		flagSet.String("format", "json", "")

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.DeleteJobs(ctx))

		var response handler.DeleteJobsCommandResponse
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, "default", response.Tube)
		require.Equal(t, 3, response.Count)
	})

	t.Run("success", func(t *testing.T) {
		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("tube", "default", "")
		flagSet.Int("count", 10, "")
		flagSet.String("format", "json", "")

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.DeleteJobs(ctx))

		var response handler.DeleteJobsCommandResponse
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, "default", response.Tube)
		require.Equal(t, 2, response.Count)
	})
}

func TestHandler_Peek(t *testing.T) {
	t.Run("empty argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.Peek(ctx))
	})

	t.Run("invalid argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.Parse([]string{"test"})

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.Peek(ctx))
	})

	t.Run("success", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Peek", 1).Return(&beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("format", "yaml", "")
		flagSet.Parse([]string{"1"})

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.Peek(ctx))

		var response handler.JobResponse
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, 1, response.ID)
		require.Equal(t, "test", response.Data)
	})
}

func TestHandler_PeekReady(t *testing.T) {
	client := &mock.Client{}
	client.On("Use", "default").Return("default", nil)
	client.On("PeekReady", mock.Anything).Return(&beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	client.On("Close", mock.Anything).Return(nil)

	h := &handler.Handler{
		Client: client,
	}

	w := new(bytes.Buffer)

	flagSet := flag.NewFlagSet("test", 0)
	flagSet.String("tube", "default", "")
	flagSet.String("format", "yaml", "")

	ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

	require.NoError(t, h.PeekReady(ctx))

	var response handler.JobResponse
	require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
	require.Equal(t, 1, response.ID)
	require.Equal(t, "test", response.Data)
}

func TestHandler_PeekDelayed(t *testing.T) {
	client := &mock.Client{}
	client.On("Use", "default").Return("default", nil)
	client.On("PeekDelayed", mock.Anything).Return(&beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	client.On("Close", mock.Anything).Return(nil)

	h := &handler.Handler{
		Client: client,
	}

	w := new(bytes.Buffer)

	flagSet := flag.NewFlagSet("test", 0)
	flagSet.String("tube", "default", "")
	flagSet.String("format", "yaml", "")

	ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

	require.NoError(t, h.PeekDelayed(ctx))

	var response handler.JobResponse
	require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
	require.Equal(t, 1, response.ID)
	require.Equal(t, "test", response.Data)
}

func TestHandler_PeekBuried(t *testing.T) {
	client := &mock.Client{}
	client.On("Use", "default").Return("default", nil)
	client.On("PeekBuried", mock.Anything).Return(&beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	client.On("Close", mock.Anything).Return(nil)

	h := &handler.Handler{
		Client: client,
	}

	w := new(bytes.Buffer)

	flagSet := flag.NewFlagSet("test", 0)
	flagSet.String("tube", "default", "")
	flagSet.String("format", "yaml", "")

	ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

	require.NoError(t, h.PeekBuried(ctx))

	var response handler.JobResponse
	require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
	require.Equal(t, 1, response.ID)
	require.Equal(t, "test", response.Data)
}

func TestHandler_Kick(t *testing.T) {
	t.Run("empty argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.Kick(ctx))
	})

	t.Run("invalid argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.Parse([]string{"test"})

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.Kick(ctx))
	})

	t.Run("success", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Use", "default").Return("default", nil)
		client.On("Kick", 3).Return(3, nil)
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("tube", "default", "")
		flagSet.String("format", "yaml", "")
		flagSet.Parse([]string{"3"})

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.Kick(ctx))

		var response handler.KickCommandResponse
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, "default", response.Tube)
		require.Equal(t, 3, response.Count)
	})
}

func TestHandler_KickJob(t *testing.T) {
	t.Run("empty argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.KickJob(ctx))
	})

	t.Run("invalid argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.Parse([]string{"test"})

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.KickJob(ctx))
	})

	t.Run("success", func(t *testing.T) {
		client := &mock.Client{}
		client.On("KickJob", 1).Return(nil)
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.Parse([]string{"1"})

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.KickJob(ctx))
	})
}

func TestHandler_KickJobs(t *testing.T) {
	client := &mock.Client{}
	client.On("Use", "default").Return("default", nil)
	client.On("PeekBuried").Times(5).Return(&beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	client.On("PeekBuried").Return(nil, beanstalk.ErrNotFound)
	client.On("KickJob", 1).Return(nil)
	client.On("Close", mock.Anything).Return(nil)

	h := &handler.Handler{
		Client: client,
	}

	t.Run("with limit", func(t *testing.T) {
		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("tube", "default", "")
		flagSet.Int("count", 3, "")
		flagSet.String("format", "json", "")

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.KickJobs(ctx))

		var response handler.KickJobsCommandResponse
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, "default", response.Tube)
		require.Equal(t, 3, response.Count)
	})

	t.Run("success", func(t *testing.T) {
		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("tube", "default", "")
		flagSet.Int("count", 10, "")
		flagSet.String("format", "json", "")

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.KickJobs(ctx))

		var response handler.KickJobsCommandResponse
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, "default", response.Tube)
		require.Equal(t, 2, response.Count)
	})
}

func TestHandler_StatsJob(t *testing.T) {
	t.Run("empty argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.StatsJob(ctx))
	})

	t.Run("invalid argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.Parse([]string{"test"})

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.StatsJob(ctx))
	})

	t.Run("success", func(t *testing.T) {
		expectedStats := &beanstalk.StatsJob{
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

		client := &mock.Client{}
		client.On("StatsJob", 1).Return(expectedStats, nil)
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("format", "yaml", "")
		flagSet.Parse([]string{"1"})

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.StatsJob(ctx))

		var response beanstalk.StatsJob
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, expectedStats.Tube, response.Tube)
		require.Equal(t, expectedStats.State, response.State)
		require.Equal(t, expectedStats.Priority, response.Priority)
		require.Equal(t, expectedStats.Age, response.Age)
		require.Equal(t, expectedStats.Delay, response.Delay)
		require.Equal(t, expectedStats.TTR, response.TTR)
		require.Equal(t, expectedStats.TimeLeft, response.TimeLeft)
		require.Equal(t, expectedStats.File, response.File)
		require.Equal(t, expectedStats.Reserves, response.Reserves)
		require.Equal(t, expectedStats.Timeouts, response.Timeouts)
		require.Equal(t, expectedStats.Releases, response.Releases)
		require.Equal(t, expectedStats.Buries, response.Buries)
		require.Equal(t, expectedStats.Kicks, response.Kicks)
	})
}

func TestHandler_StatsTube(t *testing.T) {
	t.Run("empty argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.StatsTube(ctx))
	})

	t.Run("success", func(t *testing.T) {
		expectedStats := &beanstalk.StatsTube{
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

		client := &mock.Client{}
		client.On("StatsTube", "default").Return(expectedStats, nil)
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.String("format", "yaml", "")
		flagSet.Parse([]string{"default"})

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.StatsTube(ctx))

		var response beanstalk.StatsTube
		require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
		require.Equal(t, expectedStats.CurrentJobsUrgent, response.CurrentJobsUrgent)
		require.Equal(t, expectedStats.CurrentJobsReady, response.CurrentJobsReady)
		require.Equal(t, expectedStats.CurrentJobsReserved, response.CurrentJobsReserved)
		require.Equal(t, expectedStats.CurrentJobsDelayed, response.CurrentJobsDelayed)
		require.Equal(t, expectedStats.CurrentJobsBuried, response.CurrentJobsBuried)
		require.Equal(t, expectedStats.TotalJobs, response.TotalJobs)
		require.Equal(t, expectedStats.CurrentUsing, response.CurrentUsing)
		require.Equal(t, expectedStats.CurrentWaiting, response.CurrentWaiting)
		require.Equal(t, expectedStats.CurrentWatching, response.CurrentWatching)
		require.Equal(t, expectedStats.Pause, response.Pause)
		require.Equal(t, expectedStats.CmdDelete, response.CmdDelete)
		require.Equal(t, expectedStats.CmdPauseTube, response.CmdPauseTube)
		require.Equal(t, expectedStats.PauseTimeLeft, response.PauseTimeLeft)
	})
}

func TestHandler_Stats(t *testing.T) {
	expectedStats := &beanstalk.Stats{
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

	client := &mock.Client{}
	client.On("Stats", mock.Anything).Return(expectedStats, nil)
	client.On("Close", mock.Anything).Return(nil)

	h := &handler.Handler{
		Client: client,
	}

	w := new(bytes.Buffer)

	flagSet := flag.NewFlagSet("test", 0)
	flagSet.String("format", "yaml", "")

	ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

	require.NoError(t, h.Stats(ctx))

	var response beanstalk.Stats
	require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&response))
	require.Equal(t, expectedStats.CurrentJobsUrgent, response.CurrentJobsUrgent)
	require.Equal(t, expectedStats.CurrentJobsReady, response.CurrentJobsReady)
	require.Equal(t, expectedStats.CurrentJobsReserved, response.CurrentJobsReserved)
	require.Equal(t, expectedStats.CurrentJobsDelayed, response.CurrentJobsDelayed)
	require.Equal(t, expectedStats.CurrentJobsBuried, response.CurrentJobsBuried)
	require.Equal(t, expectedStats.CmdPut, response.CmdPut)
	require.Equal(t, expectedStats.CmdPeek, response.CmdPeek)
	require.Equal(t, expectedStats.CmdPeekReady, response.CmdPeekReady)
	require.Equal(t, expectedStats.CmdPeekDelayed, response.CmdPeekDelayed)
	require.Equal(t, expectedStats.CmdPeekBuried, response.CmdPeekBuried)
	require.Equal(t, expectedStats.CmdReserve, response.CmdReserve)
	require.Equal(t, expectedStats.CmdUse, response.CmdUse)
	require.Equal(t, expectedStats.CmdWatch, response.CmdWatch)
	require.Equal(t, expectedStats.CmdIgnore, response.CmdIgnore)
	require.Equal(t, expectedStats.CmdDelete, response.CmdDelete)
	require.Equal(t, expectedStats.CmdRelease, response.CmdRelease)
	require.Equal(t, expectedStats.CmdBury, response.CmdBury)
	require.Equal(t, expectedStats.CmdKick, response.CmdKick)
	require.Equal(t, expectedStats.CmdTouch, response.CmdTouch)
	require.Equal(t, expectedStats.CmdStats, response.CmdStats)
	require.Equal(t, expectedStats.CmdStatsJob, response.CmdStatsJob)
	require.Equal(t, expectedStats.CmdStatsTube, response.CmdStatsTube)
	require.Equal(t, expectedStats.CmdListTubes, response.CmdListTubes)
	require.Equal(t, expectedStats.CmdListTubeUsed, response.CmdListTubeUsed)
	require.Equal(t, expectedStats.CmdListTubesWatched, response.CmdListTubesWatched)
	require.Equal(t, expectedStats.CmdPauseTube, response.CmdPauseTube)
	require.Equal(t, expectedStats.JobTimeouts, response.JobTimeouts)
	require.Equal(t, expectedStats.TotalJobs, response.TotalJobs)
	require.Equal(t, expectedStats.MaxJobSize, response.MaxJobSize)
	require.Equal(t, expectedStats.CurrentTubes, response.CurrentTubes)
	require.Equal(t, expectedStats.CurrentConnections, response.CurrentConnections)
	require.Equal(t, expectedStats.CurrentProducers, response.CurrentProducers)
	require.Equal(t, expectedStats.CurrentWorkers, response.CurrentWorkers)
	require.Equal(t, expectedStats.CurrentWaiting, response.CurrentWaiting)
	require.Equal(t, expectedStats.TotalConnections, response.TotalConnections)
	require.Equal(t, expectedStats.PID, response.PID)
	require.Equal(t, expectedStats.Version, response.Version)
	require.Equal(t, expectedStats.RUsageUTime, response.RUsageUTime)
	require.Equal(t, expectedStats.RUsageSTime, response.RUsageSTime)
	require.Equal(t, expectedStats.Uptime, response.Uptime)
	require.Equal(t, expectedStats.BinlogOldestIndex, response.BinlogOldestIndex)
	require.Equal(t, expectedStats.BinlogCurrentIndex, response.BinlogCurrentIndex)
	require.Equal(t, expectedStats.BinlogRecordsMigrated, response.BinlogRecordsMigrated)
	require.Equal(t, expectedStats.BinlogRecordsWritten, response.BinlogRecordsWritten)
	require.Equal(t, expectedStats.BinlogMaxSize, response.BinlogMaxSize)
	require.Equal(t, expectedStats.ID, response.ID)
	require.Equal(t, expectedStats.Hostname, response.Hostname)
}

func TestHandler_ListTubes(t *testing.T) {
	expectedTubes := []string{"default", "test"}

	client := &mock.Client{}
	client.On("ListTubes", mock.Anything).Return(expectedTubes, nil)
	client.On("Close", mock.Anything).Return(nil)

	h := &handler.Handler{
		Client: client,
	}

	w := new(bytes.Buffer)

	flagSet := flag.NewFlagSet("test", 0)
	flagSet.String("format", "yaml", "")

	ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

	require.NoError(t, h.ListTubes(ctx))

	var tubes []string
	require.NoError(t, yaml.NewDecoder(strings.NewReader(w.String())).Decode(&tubes))
	require.ElementsMatch(t, expectedTubes, tubes)
}

func TestHandler_PauseTube(t *testing.T) {
	t.Run("empty argument", func(t *testing.T) {
		client := &mock.Client{}
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		flagSet := flag.NewFlagSet("test", 0)

		ctx := cli.NewContext(nil, flagSet, nil)

		require.Error(t, h.PauseTube(ctx))
	})

	t.Run("success", func(t *testing.T) {
		client := &mock.Client{}
		client.On("PauseTube", "default", 5*time.Second).Return(nil)
		client.On("Close", mock.Anything).Return(nil)

		h := &handler.Handler{
			Client: client,
		}

		w := new(bytes.Buffer)

		flagSet := flag.NewFlagSet("test", 0)
		flagSet.Duration("delay", 5*time.Second, "")
		flagSet.Parse([]string{"default"})

		ctx := cli.NewContext(&cli.App{Writer: w}, flagSet, nil)

		require.NoError(t, h.PauseTube(ctx))
	})
}
