package middleware

import (
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/gorilla/mux"
	"net/http"
)

func Auth(provider *security.UserProvider, expectedScopes []security.Scope) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()

			if !ok {
				writer.JSON(w, http.StatusUnauthorized, response.Unauthorized())

				return
			}

			user := provider.Get(username)

			if user == nil {
				writer.JSON(w, http.StatusUnauthorized, response.Unauthorized())

				return
			}

			if !security.VerifyPassword(user.HashedPassword(), []byte(password)) {
				writer.JSON(w, http.StatusUnauthorized, response.Unauthorized())

				return
			}

			if scopes := security.VerifyScopes(user.Scopes(), expectedScopes); len(scopes) > 0 {
				data := map[string]interface{}{
					"errors": []string{
						fmt.Sprintf("required scopes %v", scopes),
					},
				}

				writer.JSON(w, http.StatusForbidden, response.Forbidden(data))

				return
			}

			ctx := security.WithAuthenticatedUser(r.Context(), user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
