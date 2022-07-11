package middleware_test

import (
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
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

	manager := security.NewTokenManager(&config.Config{
		Security: config.SecurityConfig{
			Secret:   "test",
			TokenTTL: time.Minute,
		},
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusOK, response.Success(nil))
	})

	t.Run("unauthorized", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		middleware.Auth(provider, manager, []security.Scope{}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Authorization", "test")

		middleware.Auth(provider, manager, []security.Scope{}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("unknown user", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		claims := jwt.MapClaims{
			"iss": request.URL.String(),
			"sub": "qwerty",
			"exp": time.Now().Add(5 * time.Second).Unix(),
		}

		token, err := manager.Sign(claims)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		middleware.Auth(provider, manager, []security.Scope{}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		claims := jwt.MapClaims{
			"iss": request.URL.String(),
			"sub": "test",
			"exp": time.Now().Add(5 * time.Second).Unix(),
		}

		token, err := manager.Sign(claims)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		middleware.Auth(provider, manager, []security.Scope{security.ScopeWriteJobs}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusForbidden, recorder.Code)
		require.Equal(t, `{"status":"failure","message":"Forbidden","data":{"errors":["required scopes [write:jobs]"]}}`, recorder.Body.String())
	})

	t.Run("authenticated via header", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		claims := jwt.MapClaims{
			"iss": request.URL.String(),
			"sub": "test",
			"exp": time.Now().Add(5 * time.Second).Unix(),
		}

		token, err := manager.Sign(claims)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		middleware.Auth(provider, manager, []security.Scope{security.ScopeReadServer}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("authenticated via cookies", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		claims := jwt.MapClaims{
			"iss": request.URL.String(),
			"sub": "test",
			"exp": time.Now().Add(5 * time.Second).Unix(),
		}

		token, err := manager.Sign(claims)
		if err != nil {
			t.Fatal(err)
		}

		request.AddCookie(&http.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(5 * time.Second),
			MaxAge:   5,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		middleware.Auth(provider, manager, []security.Scope{security.ScopeReadServer}).Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
	})
}
