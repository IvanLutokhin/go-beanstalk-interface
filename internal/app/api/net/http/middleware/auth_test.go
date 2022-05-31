package middleware_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuth(t *testing.T) {
	user := security.NewUser(
		"test",
		[]byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC"),
		[]security.Scope{
			security.ScopeReadServer,
			security.ScopeReadTubes,
			security.ScopeReadJobs,
		},
	)

	provider := security.NewUserProvider()
	provider.Set(user.Name(), user)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusOK, response.Success(nil))
	})

	t.Run("unauthorized", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		middleware.Auth(provider, []security.Scope{}).Middleware(handler).ServeHTTP(recorder, request)

		if code := recorder.Code; http.StatusUnauthorized != code {
			t.Errorf("expected response status code '%v', but got '%v'", http.StatusUnauthorized, code)
		}
	})

	t.Run("illegal", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Authorization", "test")

		middleware.Auth(provider, []security.Scope{}).Middleware(handler).ServeHTTP(recorder, request)

		if code := recorder.Code; http.StatusUnauthorized != code {
			t.Errorf("expected response status code '%v', but got '%v'", http.StatusUnauthorized, code)
		}
	})

	t.Run("unknown user", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "test")

		middleware.Auth(provider, []security.Scope{}).Middleware(handler).ServeHTTP(recorder, request)

		if code := recorder.Code; http.StatusUnauthorized != code {
			t.Errorf("expected response status code '%v', but got '%v'", http.StatusUnauthorized, code)
		}
	})

	t.Run("forbidden", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "password")

		middleware.Auth(provider, []security.Scope{security.ScopeWriteJobs}).Middleware(handler).ServeHTTP(recorder, request)

		if code := recorder.Code; http.StatusForbidden != code {
			t.Errorf("expected response status code '%v', but got '%v'", http.StatusForbidden, code)
		}

		if body := recorder.Body.String(); !strings.EqualFold(`{"status":"failure","message":"Forbidden","data":{"errors":["required scopes [write:jobs]"]}}`, body) {
			t.Errorf("expected error, but got '%v'", body)
		}
	})

	t.Run("authenticated", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "password")

		middleware.Auth(provider, []security.Scope{security.ScopeReadServer}).Middleware(handler).ServeHTTP(recorder, request)

		if code := recorder.Code; http.StatusOK != code {
			t.Errorf("expected response status code '%v', but got '%v'", http.StatusOK, code)
		}
	})
}
