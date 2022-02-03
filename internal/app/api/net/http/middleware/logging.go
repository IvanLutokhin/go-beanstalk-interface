package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *LoggingResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

type Logging struct {
	Logger *zap.Logger
}

func NewLogging(logger *zap.Logger) *Logging {
	return &Logging{Logger: logger}
}

func (m *Logging) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		lw.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(lw, r)

		m.Logger.
			Named("http").
			With(
				zap.Int("status_code", lw.StatusCode),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
			).
			Info("Incoming request")
	})
}
