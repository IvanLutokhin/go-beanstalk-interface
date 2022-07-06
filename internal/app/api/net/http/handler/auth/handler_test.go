package auth_test

import (
	"encoding/json"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/handler/auth"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	provider := security.NewUserProvider()
	provider.Set("test", security.NewUser(
		"test",
		[]byte("$2a$10$DwPN24dS.AL77MopVjJh/eWjwrvuRUfHLUUFTPDdwAPFLRbEzg1UC"),
		[]security.Scope{
			security.ScopeReadServer,
			security.ScopeReadTubes,
			security.ScopeReadJobs,
		},
	))

	generator := security.NewTokenGenerator(&config.Config{
		Security: config.SecurityConfig{
			Secret:   "test",
			TokenTTL: time.Second,
		},
	})

	t.Run("unknown user", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/auth/token", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("unknown", "unknown")

		auth.Token(provider, generator).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/auth/token", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "test")

		auth.Token(provider, generator).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusForbidden, recorder.Code)
	})

	t.Run("success via basic auth", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/auth/token", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.SetBasicAuth("test", "password")

		auth.Token(provider, generator).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)

		var r response.Response
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &r))

		token, ok := r.Data.(map[string]interface{})["token"].(string)
		require.True(t, ok)
		require.NotNil(t, token)
	})

	t.Run("success via body request", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPost, "/auth/token", strings.NewReader(`{"username": "test", "password": "password"}`))
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Content-Type", "application/json")

		auth.Token(provider, generator).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)

		var r response.Response
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &r))

		token, ok := r.Data.(map[string]interface{})["token"].(string)
		require.True(t, ok)
		require.NotNil(t, token)
	})
}
