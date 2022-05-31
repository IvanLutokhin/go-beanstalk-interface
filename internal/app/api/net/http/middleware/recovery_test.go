package middleware_test

import (
	"errors"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery_Middleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("test"))
	})

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	middleware.NewRecovery(zap.NewNop()).Middleware(handler).ServeHTTP(recorder, request)

	if code := recorder.Code; http.StatusInternalServerError != code {
		t.Errorf("expected response status code '%v', but got '%v'", http.StatusInternalServerError, code)
	}
}
