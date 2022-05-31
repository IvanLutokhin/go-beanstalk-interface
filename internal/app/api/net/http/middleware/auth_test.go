package middleware_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
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

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("illegal", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Authorization", "test")

		middleware.Auth(provider, []security.Scope{}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("unknown user", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "test")

		middleware.Auth(provider, []security.Scope{}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "password")

		middleware.Auth(provider, []security.Scope{security.ScopeWriteJobs}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusForbidden, recorder.Code)
		require.Equal(t, `{"status":"failure","message":"Forbidden","data":{"errors":["required scopes [write:jobs]"]}}`, recorder.Body.String())
	})

	t.Run("authenticated", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "password")

		middleware.Auth(provider, []security.Scope{security.ScopeReadServer}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
	})
}
