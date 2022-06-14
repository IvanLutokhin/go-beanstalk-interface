package middleware_test

import (
	"errors"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/stretchr/testify/require"
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

	middleware.NewRecovery().Middleware(handler).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusInternalServerError, recorder.Code)
}
