package graphql

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/pkg/beanstalk/mock"
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

	Handler(&mock.Pool{Client: &mock.Client{}}).ServeHTTP(recorder, request)

	if code := recorder.Code; 200 != code {
		t.Errorf("expected response status code '200', but got '%v'", code)
	}
}
