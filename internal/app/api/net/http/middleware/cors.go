package middleware

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"net/http"
	"regexp"
	"strings"
)

type Cors struct {
	AllowedOrigins   []*regexp.Regexp
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

func NewCors(config *config.Config) *Cors {
	c := &Cors{
		AllowCredentials: config.Http.Cors.AllowCredentials,
	}

	// Origins
	allowedOrigins := config.Http.Cors.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	for _, allowedOrigin := range allowedOrigins {
		normalizedOrigin := strings.ToLower(allowedOrigin)

		pattern := regexp.QuoteMeta(normalizedOrigin)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)

		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		c.AllowedOrigins = append(c.AllowedOrigins, re)
	}

	// Methods
	allowedMethods := config.Http.Cors.AllowedMethods
	if len(allowedMethods) == 0 {
		allowedMethods = []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		}
	}

	for _, allowedMethod := range allowedMethods {
		c.AllowedMethods = append(c.AllowedMethods, strings.ToUpper(allowedMethod))
	}

	// Headers
	allowedHeaders := config.Http.Cors.AllowedHeaders
	if len(allowedHeaders) == 0 {
		allowedHeaders = []string{"Accept", "Authorization", "Content-Type", "Origin", "X-Requested-With", "X-Auth-Token"}
	}

	for _, allowedHeader := range allowedHeaders {
		c.AllowedHeaders = append(c.AllowedHeaders, http.CanonicalHeaderKey(allowedHeader))
	}

	return c
}

func (m *Cors) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPreFlightRequest(r) {
			m.handlePreFlightRequest(w, r)

			return
		}

		m.handleRequest(w, r)

		next.ServeHTTP(w, r)
	})
}

func (m *Cors) handlePreFlightRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Origin")
	w.Header().Add("Vary", "Access-Control-Request-Method")
	w.Header().Add("Vary", "Access-Control-Request-Headers")

	origin := r.Header.Get("Origin")
	if origin == "" || !m.isAllowedOrigin(origin) {
		return
	}

	method := r.Header.Get("Access-Control-Request-Method")
	if !m.isAllowedMethod(method) {
		return
	}

	headers := strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",")
	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}

	if !m.isAllowedHeaders(headers) {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", method)
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))

	if m.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	w.WriteHeader(http.StatusOK)
}

func (m *Cors) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Vary", "Origin")

	origin := r.Header.Get("Origin")
	if origin == "" || !m.isAllowedOrigin(origin) {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)

	if m.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

func (m *Cors) isAllowedOrigin(origin string) bool {
	normalizedOrigin := strings.ToLower(origin)

	for _, allowedOrigin := range m.AllowedOrigins {
		if match := allowedOrigin.MatchString(normalizedOrigin); match {
			return true
		}
	}

	return false
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
		normalizedHeader := http.CanonicalHeaderKey(header)

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

func isPreFlightRequest(r *http.Request) bool {
	return r.Method == http.MethodOptions &&
		(r.Header.Get("Access-Control-Request-Method") != "" || r.Header.Get("Access-Control-Request-Headers") != "")
}
