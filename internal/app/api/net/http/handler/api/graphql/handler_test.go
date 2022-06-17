package graphql_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/handler/api/graphql"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/graphql", nil)
	if err != nil {
		t.Fatal(err)
	}

	graphql.Handler(&mock.Pool{}).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}
