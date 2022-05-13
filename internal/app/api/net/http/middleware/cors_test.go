package middleware

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCors_Middleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusOK, response.Success(nil))
	})

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := config.Config{
		Http: config.HttpConfig{
			Cors: config.CorsConfig{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"GET", "POST"},
				AllowHeaders:     []string{"*"},
				AllowCredentials: false,
			},
		},
	}

	NewCors(&c).Middleware(handler).ServeHTTP(recorder, request)

	if code := recorder.Code; http.StatusOK != code {
		t.Errorf("expected response status code '%v', but got '%v'", http.StatusOK, code)
	}

	if h := recorder.Header().Get("Access-Control-Allow-Origin"); h != "*" {
		t.Errorf("expected header 'Access-Control-Allow-Origin' equals '*', but got '%s'", h)
	}

	if h := recorder.Header().Get("Access-Control-Allow-Methods"); h != "GET,POST" {
		t.Errorf("expected header 'Access-Control-Allow-Origin' equals 'GET,POST', but got '%s'", h)
	}

	if h := recorder.Header().Get("Access-Control-Allow-Headers"); h != "*" {
		t.Errorf("expected header 'Access-Control-Allow-Origin' equals '*', but got '%s'", h)
	}

	if h := recorder.Header().Get("Access-Control-Allow-Credentials"); h != "false" {
		t.Errorf("expected header 'Access-Control-Allow-Origin' equals 'false', but got '%s'", h)
	}
}
