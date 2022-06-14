package middleware

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const fieldsKey = "middleware:logging:fields"

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

		ctx := context.WithValue(r.Context(), fieldsKey, map[string]interface{}{})

		next.ServeHTTP(lw, r.WithContext(ctx))

		fields := make([]zap.Field, 0)
		if data, ok := ctx.Value(fieldsKey).(map[string]interface{}); ok {
			for k, v := range data {
				fields = append(fields, zap.Any(k, v))
			}
		}

		m.Logger.
			Named("http").
			With(
				zap.Int("status_code", lw.StatusCode),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
			).
			With(fields...).
			Info("Incoming request")
	})
}
