package middleware

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"net/http"
	"strconv"
	"strings"
)

type Cors struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

func NewCors(config *config.Config) *Cors {
	return &Cors{
		AllowOrigins:     config.Http.Cors.AllowOrigins,
		AllowMethods:     config.Http.Cors.AllowMethods,
		AllowHeaders:     config.Http.Cors.AllowHeaders,
		AllowCredentials: config.Http.Cors.AllowCredentials,
	}
}

func (m *Cors) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(m.AllowOrigins, ","))
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.AllowMethods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.AllowHeaders, ","))
		w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(m.AllowCredentials))

		next.ServeHTTP(w, r)
	})
}
