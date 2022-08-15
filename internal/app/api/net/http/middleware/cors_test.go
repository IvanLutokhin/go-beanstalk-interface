package middleware_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors_Middleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusOK, response.Success(nil))
	})

	t.Run("preflight / success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodOptions, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Origin", "http://127.0.0.1")
		request.Header.Set("Access-Control-Request-Method", http.MethodGet)
		request.Header.Set("Access-Control-Request-Headers", "Authorization")

		middleware.NewCors().Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.ElementsMatch(t, []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"}, recorder.Header().Values("Vary"))
		require.Equal(t, "http://127.0.0.1", recorder.Header().Get("Access-Control-Allow-Origin"))
		require.Equal(t, http.MethodGet, recorder.Header().Get("Access-Control-Allow-Methods"))
		require.Equal(t, "Authorization", recorder.Header().Get("Access-Control-Allow-Headers"))
		require.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("preflight / empty origin", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodOptions, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Access-Control-Request-Method", http.MethodGet)
		request.Header.Set("Access-Control-Request-Headers", "Authorization")

		middleware.NewCors().Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.ElementsMatch(t, []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"}, recorder.Header().Values("Vary"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Methods"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Headers"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("preflight / not allowed method", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodOptions, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Origin", "http://127.0.0.1")
		request.Header.Set("Access-Control-Request-Method", "TEST")
		request.Header.Set("Access-Control-Request-Headers", "Authorization")

		middleware.NewCors().Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.ElementsMatch(t, []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"}, recorder.Header().Values("Vary"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Methods"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Headers"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("preflight / not allowed headers", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodOptions, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Origin", "http://127.0.0.1")
		request.Header.Set("Access-Control-Request-Method", http.MethodGet)
		request.Header.Set("Access-Control-Request-Headers", "X-Test")

		middleware.NewCors().Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.ElementsMatch(t, []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"}, recorder.Header().Values("Vary"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Methods"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Headers"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("request / success", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		request.Header.Set("Origin", "http://127.0.0.1")
		request.Header.Set("Authorization", "test")

		middleware.NewCors().Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.ElementsMatch(t, []string{"Origin"}, recorder.Header().Values("Vary"))
		require.Equal(t, "http://127.0.0.1", recorder.Header().Get("Access-Control-Allow-Origin"))
		require.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("request / empty origin", func(t *testing.T) {
		recorder := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api", nil)
		if err != nil {
			t.Fatal(err)
		}

		middleware.NewCors().Middleware(handler).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.ElementsMatch(t, []string{"Origin"}, recorder.Header().Values("Vary"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Origin"))
		require.Empty(t, recorder.Header().Get("Access-Control-Allow-Credentials"))
	})
}
