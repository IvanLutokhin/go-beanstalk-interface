package middleware_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging_Middleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusAccepted, response.Success(nil))
	})

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	middleware.NewLogging(zap.NewNop()).Middleware(handler).ServeHTTP(recorder, request)

	if code := recorder.Code; http.StatusAccepted != code {
		t.Errorf("expected response status code '%v', but got '%v'", http.StatusAccepted, code)
	}
}
