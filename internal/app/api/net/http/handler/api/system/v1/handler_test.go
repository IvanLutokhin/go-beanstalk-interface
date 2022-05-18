package v1

import (
	"encoding/json"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/api"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/pkg/beanstalk/mock"
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/embed"
	"github.com/gorilla/mux"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"
)

func TestGetEmbedFiles(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/system/v1/swagger.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.StripPrefix("/api/system/v1", GetEmbedFiles(http.FS(embed.FSFunc(func(name string) (fs.File, error) {
		return api.SystemV1EmbedFS.Open(path.Join("system/v1", name))
	})))).ServeHTTP(recorder, request)

	if code := recorder.Code; http.StatusOK != code {
		t.Errorf("expected response status code '%v', but got '%v'", http.StatusOK, code)
	}
}

func TestGetServerStats(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/system/v1/server/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	beanstalk.NewHTTPHandlerAdapter(pool, GetServerStats()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)

	body, err := UnmarshalBody(recorder)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := body.Data.(map[string]interface{})["stats"]; !ok {
		t.Error("unexpected data in response")
	}
}

func TestGetTubes(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/system/v1/tubes", nil)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	beanstalk.NewHTTPHandlerAdapter(pool, GetTubes()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)

	body, err := UnmarshalBody(recorder)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := body.Data.(map[string]interface{})["tubes"]; !ok {
		t.Error("unexpected data in response")
	}
}

func TestGetTubeStats(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("tube stats", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/tubes/default/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"name": "default",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, GetTubeStats()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)

		body, err := UnmarshalBody(recorder)
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := body.Data.(map[string]interface{})["stats"]; !ok {
			t.Error("unexpected data in response")
		}
	})

	t.Run("tube stats / not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/tubes/not_found/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"name": "not_found",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, GetTubeStats()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestCreateJob(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("create job", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(`{"tube": "default", "data": "test"}`))
		if err != nil {
			t.Fatal(err)
		}

		beanstalk.NewHTTPHandlerAdapter(pool, CreateJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusCreated)

		body, err := UnmarshalBody(recorder)
		if err != nil {
			t.Fatal(err)
		}

		tube := body.Data.(map[string]interface{})["tube"]
		if !strings.EqualFold("default", tube.(string)) {
			t.Errorf("excpeted tube 'default', but got '%s'", tube)
		}

		id := body.Data.(map[string]interface{})["id"]
		if 1 != id.(float64) {
			t.Errorf("excpeted job id '1', but got '%d'", id)
		}
	})

	t.Run("create job / bad JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs", strings.NewReader("test"))
		if err != nil {
			t.Fatal(err)
		}

		beanstalk.NewHTTPHandlerAdapter(pool, CreateJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusBadRequest)
	})
}

func TestGetJob(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get job", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/jobs/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, GetJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)

		body, err := UnmarshalBody(recorder)
		if err != nil {
			t.Fatal(err)
		}

		data, ok := body.Data.(map[string]interface{})["data"]
		if !ok {
			t.Error("unexpected data in response")
		}

		if !strings.EqualFold("test", data.(string)) {
			t.Errorf("excpeted data 'test', but got '%s'", data)
		}
	})

	t.Run("get job / not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/jobs/999", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, GetJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestBuryJob(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("bury job", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/bury", strings.NewReader(`{"priority": 0}`))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, BuryJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)
	})

	t.Run("bury job / bad JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/bury", strings.NewReader("test"))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, BuryJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusBadRequest)
	})

	t.Run("bury job / not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/999/bury", strings.NewReader(`{"priority": 100}`))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, BuryJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestDeleteJob(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("delete job", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/delete", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, DeleteJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)
	})

	t.Run("delete job / not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/999/delete", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, DeleteJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestKickJob(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("kick job", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/kick", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, KickJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)
	})

	t.Run("kick job / not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/999/kick", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, KickJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestReleaseJob(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("release job", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/release", strings.NewReader(`{"priority": 0, "delay": 0}`))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, ReleaseJob()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)
	})

	t.Run("release job / bad JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/1/release", strings.NewReader("test"))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, ReleaseJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusBadRequest)
	})

	t.Run("release job / not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/api/system/v1/jobs/999/release", strings.NewReader(`{"priority": 100, "delay": 100}`))
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, ReleaseJob()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

func TestGetJobsStats(t *testing.T) {
	pool, err := beanstalk.NewPool(func() (beanstalk.Client, error) { return &mock.Client{}, nil }, 1, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get job stats", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/jobs/1/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "1",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, GetJobStats()).ServeHTTP(recorder, request)

		AssertResponseSuccess(t, recorder, http.StatusOK)

		body, err := UnmarshalBody(recorder)
		if err != nil {
			t.Fatal(err)
		}

		if _, ok := body.Data.(map[string]interface{})["stats"]; !ok {
			t.Error("unexpected data in response")
		}
	})

	t.Run("get job stats / not found", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/system/v1/jobs/999/stats", nil)
		if err != nil {
			t.Fatal(err)
		}

		request = mux.SetURLVars(request, map[string]string{
			"id": "999",
		})

		beanstalk.NewHTTPHandlerAdapter(pool, GetJobStats()).ServeHTTP(recorder, request)

		AssertResponseFailure(t, recorder, http.StatusNotFound)
	})
}

// Helpers

func UnmarshalBody(recorder *httptest.ResponseRecorder) (response.Response, error) {
	var body response.Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		return body, err
	}

	return body, nil
}

func AssertResponse(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode int, expectedStatus string) {
	if code := recorder.Code; expectedCode != code {
		t.Errorf("expected response status code '%v', but got '%v'", expectedCode, code)
	}

	body, err := UnmarshalBody(recorder)
	if err != nil {
		t.Fatal(err)
	}

	if expectedStatus != body.Status {
		t.Errorf("expected response status '%s', but got '%s'", expectedStatus, body.Status)
	}
}

func AssertResponseSuccess(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode int) {
	AssertResponse(t, recorder, expectedCode, response.StatusSuccess)
}

func AssertResponseFailure(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode int) {
	AssertResponse(t, recorder, expectedCode, response.StatusFailure)
}
