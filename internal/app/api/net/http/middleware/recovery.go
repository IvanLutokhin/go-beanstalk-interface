package middleware

import (
	"errors"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"net/http"
)

type Recovery struct{}

func NewRecovery() *Recovery {
	return &Recovery{}
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

				if fields, ok := r.Context().Value("middleware:logging:fields").(map[string]interface{}); ok {
					fields["error"] = err
					fields["panic_occurred"] = true
				}

				writer.JSON(w, http.StatusInternalServerError, response.InternalServerError())
			}
		}()

		next.ServeHTTP(w, r)
	})
}
