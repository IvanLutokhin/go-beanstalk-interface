package auth

import (
	"encoding/json"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

func Token(config *config.Config, provider *security.UserProvider, manager *security.TokenManager) http.Handler {
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

		claims := jwt.MapClaims{
			"iss": r.URL.String(),
			"sub": username,
			"exp": time.Now().Add(config.Security.TokenTTL).Unix(),
		}

		token, err := manager.Sign(claims)
		if err != nil {
			panic(err)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(config.Security.TokenTTL),
			MaxAge:   int(config.Security.TokenTTL.Seconds()),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "logged_in",
			Value:    "1",
			Path:     "/",
			Expires:  time.Now().Add(config.Security.TokenTTL),
			MaxAge:   int(config.Security.TokenTTL.Seconds()),
			SameSite: http.SameSiteLaxMode,
		})

		data := map[string]interface{}{
			"access_token": token,
			"token_type":   "bearer",
			"expires_in":   config.Security.TokenTTL.Seconds(),
		}

		writer.JSON(w, http.StatusOK, response.Success(data))
	})
}

func Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "logged_in",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			SameSite: http.SameSiteLaxMode,
		})

		writer.JSON(w, http.StatusOK, response.Success(nil))
	})
}
