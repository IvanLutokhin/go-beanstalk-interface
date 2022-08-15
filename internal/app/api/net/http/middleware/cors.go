package middleware

import (
	"net/http"
	"strings"
)

type Cors struct {
	AllowedMethods []string
	AllowedHeaders []string
}

func NewCors() *Cors {
	return &Cors{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"Origin",
			"X-Requested-With",
		},
	}
}

func (m *Cors) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPreflightRequest(r) {
			m.handlePreflightRequest(w, r)
		} else {
			m.handleRequest(w, r)

			next.ServeHTTP(w, r)
		}
	})
}

func isPreflightRequest(r *http.Request) bool {
	return r.Method == http.MethodOptions &&
		(r.Header.Get("Access-Control-Request-Method") != "" || r.Header.Get("Access-Control-Request-Headers") != "")
}

func (m *Cors) handlePreflightRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")

	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	method := r.Header.Get("Access-Control-Request-Method")
	if !m.isAllowedMethod(method) {
		return
	}

	headers := strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",")
	if !m.isAllowedHeaders(headers) {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", method)
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusOK)
}

func (m *Cors) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Origin")

	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func (m *Cors) isAllowedMethod(method string) bool {
	normalizedMethod := strings.ToUpper(method)

	for _, allowedMethod := range m.AllowedMethods {
		if allowedMethod == normalizedMethod {
			return true
		}
	}

	return false
}

func (m *Cors) isAllowedHeaders(headers []string) bool {
	for _, header := range headers {
		normalizedHeader := http.CanonicalHeaderKey(strings.TrimSpace(header))

		found := false
		for _, allowedHeader := range m.AllowedHeaders {
			if allowedHeader == normalizedHeader {
				found = true

				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
