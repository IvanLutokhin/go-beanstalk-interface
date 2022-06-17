package graphql_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/handler/api/graphql"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodPost, "/api/graphql", strings.NewReader("{ \"query\": \"{ __typename }\" }"))
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")

	graphql.Handler(&mock.Pool{}).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)
}
