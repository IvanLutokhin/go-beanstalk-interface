package v1_test

import (
	"encoding/json"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/api"
	v1 "github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/handler/api/system/v1"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/embed"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
	"time"
)

func TestGetEmbedFiles(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/system/v1/swagger.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.StripPrefix("/api/system/v1", v1.GetEmbedFiles(http.FS(embed.FSFunc(func(name string) (fs.File, error) {
		return api.SystemV1EmbedFS.Open(path.Join("system/v1", name))
	})))).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetServerStats(t *testing.T) {
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

	client := &mock.Client{}
	client.On("Stats").Return(expectedStats, nil)

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/system/v1/server/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	beanstalk.NewHTTPHandlerAdapter(mock.NewPool(client), v1.GetServerStats()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)

	body, err := UnmarshalBody(recorder)
	if err != nil {
		t.Fatal(err)
	}

	stats, ok := body.Data.(map[string]interface{})["stats"].(map[string]interface{})

	require.True(t, ok)
	require.Equal(t, float64(expectedStats.CurrentJobsUrgent), stats["currentJobsUrgent"])
	require.Equal(t, float64(expectedStats.CurrentJobsReady), stats["currentJobsReady"])
	require.Equal(t, float64(expectedStats.CurrentJobsReserved), stats["currentJobsReserved"])
	require.Equal(t, float64(expectedStats.CurrentJobsDelayed), stats["currentJobsDelayed"])
	require.Equal(t, float64(expectedStats.CurrentJobsBuried), stats["currentJobsBuried"])
	require.Equal(t, float64(expectedStats.CmdPut), stats["cmdPut"])
	require.Equal(t, float64(expectedStats.CmdPeek), stats["cmdPeek"])
	require.Equal(t, float64(expectedStats.CmdPeekReady), stats["cmdPeekReady"])
	require.Equal(t, float64(expectedStats.CmdPeekDelayed), stats["cmdPeekDelayed"])
	require.Equal(t, float64(expectedStats.CmdPeekBuried), stats["cmdPeekBuried"])
	require.Equal(t, float64(expectedStats.CmdReserve), stats["cmdReserve"])
	require.Equal(t, float64(expectedStats.CmdUse), stats["cmdUse"])
	require.Equal(t, float64(expectedStats.CmdWatch), stats["cmdWatch"])
	require.Equal(t, float64(expectedStats.CmdIgnore), stats["cmdIgnore"])
	require.Equal(t, float64(expectedStats.CmdDelete), stats["cmdDelete"])
	require.Equal(t, float64(expectedStats.CmdRelease), stats["cmdRelease"])
	require.Equal(t, float64(expectedStats.CmdBury), stats["cmdBury"])
	require.Equal(t, float64(expectedStats.CmdKick), stats["cmdKick"])
	require.Equal(t, float64(expectedStats.CmdTouch), stats["cmdTouch"])
	require.Equal(t, float64(expectedStats.CmdStats), stats["cmdStats"])
	require.Equal(t, float64(expectedStats.CmdStatsJob), stats["cmdStatsJob"])
	require.Equal(t, float64(expectedStats.CmdStatsTube), stats["cmdStatsTube"])
	require.Equal(t, float64(expectedStats.CmdListTubes), stats["cmdListTubes"])
	require.Equal(t, float64(expectedStats.CmdListTubeUsed), stats["cmdListTubeUsed"])
	require.Equal(t, float64(expectedStats.CmdListTubesWatched), stats["cmdListTubesWatched"])
	require.Equal(t, float64(expectedStats.CmdPauseTube), stats["cmdPauseTube"])
	require.Equal(t, float64(expectedStats.JobTimeouts), stats["jobTimeouts"])
	require.Equal(t, float64(expectedStats.TotalJobs), stats["totalJobs"])
	require.Equal(t, float64(expectedStats.MaxJobSize), stats["maxJobSize"])
	require.Equal(t, float64(expectedStats.CurrentTubes), stats["currentTubes"])
	require.Equal(t, float64(expectedStats.CurrentConnections), stats["currentConnections"])
	require.Equal(t, float64(expectedStats.CurrentProducers), stats["currentProducers"])
	require.Equal(t, float64(expectedStats.CurrentWorkers), stats["currentWorkers"])
	require.Equal(t, float64(expectedStats.CurrentWaiting), stats["currentWaiting"])
	require.Equal(t, float64(expectedStats.TotalConnections), stats["totalConnections"])
	require.Equal(t, float64(expectedStats.PID), stats["pid"])
	require.Equal(t, expectedStats.Version, stats["version"])
	require.Equal(t, expectedStats.RUsageUTime, stats["rUsageUTime"])
	require.Equal(t, expectedStats.RUsageSTime, stats["rUsageSTime"])
	require.Equal(t, float64(expectedStats.Uptime), stats["uptime"])
	require.Equal(t, float64(expectedStats.BinlogOldestIndex), stats["binlogOldestIndex"])
	require.Equal(t, float64(expectedStats.BinlogCurrentIndex), stats["binlogCurrentIndex"])
	require.Equal(t, float64(expectedStats.BinlogRecordsMigrated), stats["binlogRecordsMigrated"])
	require.Equal(t, float64(expectedStats.BinlogRecordsWritten), stats["binlogRecordsWritten"])
	require.Equal(t, float64(expectedStats.BinlogMaxSize), stats["binlogMaxSize"])
	require.Equal(t, expectedStats.ID, stats["id"])
	require.Equal(t, expectedStats.Hostname, stats["hostname"])
}

func TestGetTubes(t *testing.T) {
	expectedTubes := []string{"default", "test"}

	client := &mock.Client{}
	client.On("ListTubes").Return(expectedTubes, nil)

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/system/v1/tubes", nil)
	if err != nil {
		t.Fatal(err)
	}

	beanstalk.NewHTTPHandlerAdapter(mock.NewPool(client), v1.GetTubes()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)

	body, err := UnmarshalBody(recorder)
	if err != nil {
		t.Fatal(err)
	}

	tubes, ok := body.Data.(map[string]interface{})["tubes"]

	require.True(t, ok)
	require.ElementsMatch(t, expectedTubes, tubes)
}

func TestGetTubeStats(t *testing.T) {
	expectedStats := beanstalk.StatsTube{
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
	client.On("StatsTube", "not_found").Return(beanstalk.StatsTube{}, beanstalk.ErrNotFound)

	pool := mock.NewPool(client)

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/tubes/default/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"name": "default",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.GetTubeStats()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)

		body, err := UnmarshalBody(recorder)
		if err != nil {
			t.Fatal(err)
		}

		stats, ok := body.Data.(map[string]interface{})["stats"].(map[string]interface{})

		require.True(t, ok)
		require.Equal(t, expectedStats.Name, stats["name"])
		require.Equal(t, float64(expectedStats.CurrentJobsUrgent), stats["currentJobsUrgent"])
		require.Equal(t, float64(expectedStats.CurrentJobsReady), stats["currentJobsReady"])
		require.Equal(t, float64(expectedStats.CurrentJobsReserved), stats["currentJobsReserved"])
		require.Equal(t, float64(expectedStats.CurrentJobsDelayed), stats["currentJobsDelayed"])
		require.Equal(t, float64(expectedStats.CurrentJobsBuried), stats["currentJobsBuried"])
		require.Equal(t, float64(expectedStats.TotalJobs), stats["totalJobs"])
		require.Equal(t, float64(expectedStats.CurrentUsing), stats["currentUsing"])
		require.Equal(t, float64(expectedStats.CurrentWaiting), stats["currentWaiting"])
		require.Equal(t, float64(expectedStats.CurrentWatching), stats["currentWatching"])
		require.Equal(t, float64(expectedStats.Pause), stats["pause"])
		require.Equal(t, float64(expectedStats.CmdDelete), stats["cmdDelete"])
		require.Equal(t, float64(expectedStats.CmdPauseTube), stats["cmdPauseTube"])
		require.Equal(t, float64(expectedStats.PauseTimeLeft), stats["pauseTimeLeft"])
	})

	t.Run("not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/tubes/not_found/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"name": "not_found",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.GetTubeStats()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestCreateJob(t *testing.T) {
	client := &mock.Client{}
	client.On("Use", "default").Return("default", nil)
	client.On("Put", uint32(0), time.Duration(0), time.Duration(0), []byte("test")).Return(1, nil)

	pool := mock.NewPool(client)

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(`{"tube": "default", "data": "test"}`))
		if err != nil {
			t.Fatal(err)
		}

		beanstalk.NewHTTPHandlerAdapter(pool, v1.CreateJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusCreated)

		body, err := UnmarshalBody(recorder)
		if err != nil {
			t.Fatal(err)
		}

		tube, ok := body.Data.(map[string]interface{})["tube"]

		require.True(t, ok)
		require.Equal(t, "default", tube)

		id, ok := body.Data.(map[string]interface{})["id"]

		require.True(t, ok)
		require.Equal(t, float64(1), id)
	})

	t.Run("bad JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs", strings.NewReader("test"))
		if err != nil {
			t.Fatal(err)
		}

		beanstalk.NewHTTPHandlerAdapter(pool, v1.CreateJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusBadRequest)
	})
}

func TestGetJob(t *testing.T) {
	client := &mock.Client{}
	client.On("Peek", 1).Return(beanstalk.Job{ID: 1, Data: []byte("test")}, nil)
	client.On("Peek", 999).Return(beanstalk.Job{}, beanstalk.ErrNotFound)

	pool := mock.NewPool(client)

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/jobs/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.GetJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)

		body, err := UnmarshalBody(recorder)
		if err != nil {
			t.Fatal(err)
		}

		data, ok := body.Data.(map[string]interface{})["data"]

		require.True(t, ok)
		require.Equal(t, "test", data)
	})

	t.Run("not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/jobs/999", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.GetJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestBuryJob(t *testing.T) {
	client := &mock.Client{}
	client.On("Bury", 1, uint32(0)).Return(nil)
	client.On("Bury", 999, uint32(100)).Return(beanstalk.ErrNotFound)

	pool := mock.NewPool(client)

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/bury", strings.NewReader(`{"priority": 0}`))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.BuryJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)
	})

	t.Run("bad JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/bury", strings.NewReader("test"))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.BuryJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/999/bury", strings.NewReader(`{"priority": 100}`))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.BuryJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestDeleteJob(t *testing.T) {
	client := &mock.Client{}
	client.On("Delete", 1).Return(nil)
	client.On("Delete", 999).Return(beanstalk.ErrNotFound)

	pool := mock.NewPool(client)

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/delete", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.DeleteJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)
	})

	t.Run("not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/999/delete", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.DeleteJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestKickJob(t *testing.T) {
	client := &mock.Client{}
	client.On("KickJob", 1).Return(nil)
	client.On("KickJob", 999).Return(beanstalk.ErrNotFound)

	pool := mock.NewPool(client)

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/kick", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.KickJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)
	})

	t.Run("not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/999/kick", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.KickJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestReleaseJob(t *testing.T) {
	client := &mock.Client{}
	client.On("Release", 1, uint32(0), 5*time.Second).Return(nil)
	client.On("Release", 999, uint32(100), 100*time.Second).Return(beanstalk.ErrNotFound)

	pool := mock.NewPool(client)

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/release", strings.NewReader(`{"priority": 0, "delay": "5s"}`))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.ReleaseJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)
	})

	t.Run("bad JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/release", strings.NewReader("test"))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.ReleaseJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusBadRequest)
	})

	t.Run("not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/999/release", strings.NewReader(`{"priority": 100, "delay": "100s"}`))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.ReleaseJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestGetJobsStats(t *testing.T) {
	expectedStats := beanstalk.StatsJob{
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

	client := &mock.Client{}
	client.On("StatsJob", 1).Return(expectedStats, nil)
	client.On("StatsJob", 999).Return(beanstalk.StatsJob{}, beanstalk.ErrNotFound)

	pool := mock.NewPool(client)

	t.Run("success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/jobs/1/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.GetJobStats()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)

		body, err := UnmarshalBody(recorder)
		if err != nil {
			t.Fatal(err)
		}

		stats, ok := body.Data.(map[string]interface{})["stats"].(map[string]interface{})

		require.True(t, ok)
		require.Equal(t, float64(expectedStats.ID), stats["id"])
		require.Equal(t, expectedStats.Tube, stats["tube"])
		require.Equal(t, expectedStats.State, stats["state"])
		require.Equal(t, float64(expectedStats.Priority), stats["priority"])
		require.Equal(t, float64(expectedStats.Age), stats["age"])
		require.Equal(t, float64(expectedStats.Delay), stats["delay"])
		require.Equal(t, float64(expectedStats.TTR), stats["ttr"])
		require.Equal(t, float64(expectedStats.TimeLeft), stats["timeLeft"])
		require.Equal(t, float64(expectedStats.File), stats["file"])
		require.Equal(t, float64(expectedStats.Reserves), stats["reserves"])
		require.Equal(t, float64(expectedStats.Timeouts), stats["timeouts"])
		require.Equal(t, float64(expectedStats.Releases), stats["releases"])
		require.Equal(t, float64(expectedStats.Buries), stats["buries"])
		require.Equal(t, float64(expectedStats.Kicks), stats["kicks"])
	})

	t.Run("not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/jobs/999/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, v1.GetJobStats()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

// Helpers

func UnmarshalBody(recorder *httptest.ResponseRecorder) (body response.Response, err error) {
	err = json.Unmarshal(recorder.Body.Bytes(), &body)

	return
}

func AssertResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode int, expectedStatus string) {
	require.Equal(t, expectedCode, recorder.Code)

	body, err := UnmarshalBody(recorder)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, expectedStatus, body.Status)
}

func AssertResponseSuccess(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode int) {
	AssertResponse(t, recorder, expectedCode, response.StatusSuccess)
}

func AssertResponseFailure(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode int) {
	AssertResponse(t, recorder, expectedCode, response.StatusFailure)
}
