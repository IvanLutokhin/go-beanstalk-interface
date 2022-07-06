package auth

import (
	"encoding/json"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"net/http"
)

func Token(provider *security.UserProvider, generator *security.TokenGenerator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			var request = struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}{}

			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				panic(err)
			}

			username, password = request.Username, request.Password
		}

		user := provider.Get(username)
		if user == nil || !security.VerifyPassword(user.HashedPassword(), []byte(password)) {
			data := map[string]interface{}{
				"errors": []string{
					"invalid credentials",
				},
			}

			writer.JSON(w, http.StatusForbidden, response.Forbidden(data))

			return
		}

		claims := &security.TokenClaims{
			Request: r,
			User:    user,
		}

		token, err := generator.Generate(claims)
		if err != nil {
			panic(err)
		}

		writer.JSON(w, http.StatusOK, response.Success(map[string]string{"token": token}))
	})
}
