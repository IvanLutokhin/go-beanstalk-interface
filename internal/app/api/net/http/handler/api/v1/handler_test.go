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

	request, err := http.NewRequest(http.MethodGet, "/api/v1/swagger.json", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.StripPrefix("/api/v1", GetEmbedFiles(http.FS(embed.FSFunc(func(name string) (fs.File, error) {
		return api.V1EmbedFS.Open(path.Join("v1", name))
	})))).ServeHTTP(recorder, request)

	if code := recorder.Code; http.StatusOK != code {
		t.Errorf("expected response status code '%v', but got '%v'", http.StatusOK, code)
	}
}

func TestGetServerStats(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/server/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	MockBeanstalkHandler(GetServerStats()).ServeHTTP(recorder, request)

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

	request, err := http.NewRequest(http.MethodGet, "/api/v1/tubes", nil)
	if err != nil {
		t.Fatal(err)
	}

	MockBeanstalkHandler(GetTubes()).ServeHTTP(recorder, request)

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
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/tubes/default/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"name": "default",
	})

	MockBeanstalkHandler(GetTubeStats()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)

	body, err := UnmarshalBody(recorder)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := body.Data.(map[string]interface{})["stats"]; !ok {
		t.Error("unexpected data in response")
	}
}

func TestGetTubeStatsNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/tubes/not_found/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"name": "not_found",
	})

	MockBeanstalkHandler(GetTubeStats()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusNotFound)
}

func TestCreateJob(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(`{"tube": "default", "data": "test"}`))
	if err != nil {
		t.Fatal(err)
	}

	MockBeanstalkHandler(CreateJob()).ServeHTTP(recorder, request)

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
}

func TestCreateJobBadJSON(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader("test"))
	if err != nil {
		t.Fatal(err)
	}

	MockBeanstalkHandler(CreateJob()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusBadRequest)
}

func TestGetJob(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/jobs/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "1",
	})

	MockBeanstalkHandler(GetJob()).ServeHTTP(recorder, request)

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
}

func TestGetJobNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/jobs/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "999",
	})

	MockBeanstalkHandler(GetJob()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusNotFound)
}

func TestBuryJob(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/1/bury", strings.NewReader(`{"priority": 0}`))
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "1",
	})

	MockBeanstalkHandler(BuryJob()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)
}

func TestBuryJobBadJson(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/1/bury", strings.NewReader("test"))
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "1",
	})

	MockBeanstalkHandler(BuryJob()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusBadRequest)
}

func TestBuryJobNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/999/bury", strings.NewReader(`{"priority": 100}`))
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "999",
	})

	MockBeanstalkHandler(BuryJob()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusNotFound)
}

func TestDeleteJob(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/1/delete", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "1",
	})

	MockBeanstalkHandler(DeleteJob()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)
}

func TestDeleteJobNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/999/delete", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "999",
	})

	MockBeanstalkHandler(DeleteJob()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusNotFound)
}

func TestKickJob(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/1/kick", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "1",
	})

	MockBeanstalkHandler(KickJob()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)
}

func TestKickJobNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/999/kick", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "999",
	})

	MockBeanstalkHandler(KickJob()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusNotFound)
}

func TestReleaseJob(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/1/release", strings.NewReader(`{"priority": 0, "delay": 0}`))
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "1",
	})

	MockBeanstalkHandler(ReleaseJob()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)
}

func TestReleaseJobBadJson(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/1/release", strings.NewReader("test"))
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "1",
	})

	MockBeanstalkHandler(ReleaseJob()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusBadRequest)
}

func TestReleaseJobNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/v1/jobs/999/release", strings.NewReader(`{"priority": 100, "delay": 100}`))
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "999",
	})

	MockBeanstalkHandler(ReleaseJob()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusNotFound)
}

func TestGetJobsStats(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/jobs/1/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "1",
	})

	MockBeanstalkHandler(GetJobStats()).ServeHTTP(recorder, request)

	AssertResponseSuccess(t, recorder, http.StatusOK)

	body, err := UnmarshalBody(recorder)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := body.Data.(map[string]interface{})["stats"]; !ok {
		t.Error("unexpected data in response")
	}
}

func TestGetJobStatsNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/jobs/999/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"id": "999",
	})

	MockBeanstalkHandler(GetJobStats()).ServeHTTP(recorder, request)

	AssertResponseFailure(t, recorder, http.StatusNotFound)
}

// Helpers

func MockBeanstalkHandler(handler beanstalk.Handler) http.Handler {
	return beanstalk.NewHTTPHandlerAdapter(&mock.Pool{Client: &mock.Client{}}, handler)
}

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
