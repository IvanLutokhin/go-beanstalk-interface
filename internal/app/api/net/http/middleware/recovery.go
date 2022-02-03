package middleware

import (
	"errors"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"go.uber.org/zap"
	"net/http"
)

type Recovery struct {
	Logger *zap.Logger
}

func NewRecovery(logger *zap.Logger) *Recovery {
	return &Recovery{Logger: logger}
}

func (m *Recovery) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				var err error
				switch v := e.(type) {
				case error:
					err = v
				case string:
					err = errors.New(v)
				default:
					err = errors.New("unknown error")
				}

				m.Logger.
					Named("http").
					With(zap.Error(err)).
					Error("Internal Server Error")

				writer.JSON(w, http.StatusInternalServerError, response.InternalServerError())
			}
		}()

		next.ServeHTTP(w, r)
	})
}
