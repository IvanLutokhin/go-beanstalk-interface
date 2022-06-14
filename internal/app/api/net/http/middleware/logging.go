package middleware

import (
	"bytes"
	"context"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	fieldsKey = "middleware:logging:fields"
	loggerMsg = "Incoming request"
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
		lw := &LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		var body bytes.Buffer
		if r.Body != nil {
			_, err := body.ReadFrom(r.Body)
			if err != nil {
				panic(err)
			}

			r.Body = ioutil.NopCloser(bytes.NewBuffer(body.Bytes()))
		}

		ctx := context.WithValue(r.Context(), fieldsKey, map[string]interface{}{})

		start := time.Now()

		next.ServeHTTP(lw, r.WithContext(ctx))

		fields := make([]zap.Field, 0)
		if data, ok := ctx.Value(fieldsKey).(map[string]interface{}); ok {
			for k, v := range data {
				fields = append(fields, zap.Any(k, v))
			}
		}

		logger := m.Logger.
			Named("http").
			With(
				zap.Int("status_code", lw.StatusCode),
				zap.Duration("duration", time.Since(start)),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
			).
			With(fields...)

		if lw.StatusCode < 400 {
			logger.Info(loggerMsg)
		} else if lw.StatusCode >= 400 && lw.StatusCode < 500 {
			logger.Warn(loggerMsg)
		} else {
			logger.
				With(zap.String("body", string(body.Next(1024)))).
				Error(loggerMsg)
		}
	})
}
