package middleware

import (
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-jwt/jwt/v4/request"
	"github.com/gorilla/mux"
	"net/http"
)

func Auth(provider *security.UserProvider, manager *security.TokenManager, expectedScopes []security.Scope) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := request.AuthorizationHeaderExtractor.ExtractToken(r)
			if err != nil {
				cookie, err := r.Cookie("access_token")
				if err != nil {
					writer.JSON(w, http.StatusUnauthorized, response.Unauthorized())

					return
				}

				tokenString = cookie.Value
			}

			token, err := manager.Extract(tokenString)
			if err != nil {
				writer.JSON(w, http.StatusUnauthorized, response.Unauthorized())

				return
			}

			username, ok := token.Claims.(jwt.MapClaims)["sub"].(string)
			if !ok {
				writer.JSON(w, http.StatusUnauthorized, response.Unauthorized())

				return
			}

			user := provider.Get(username)
			if user == nil {
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
