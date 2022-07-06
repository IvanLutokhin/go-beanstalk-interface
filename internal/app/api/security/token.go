package security

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type TokenClaims struct {
	Request *http.Request
	User    *User
}

type TokenGenerator struct {
	secret string
	ttl    time.Duration
}

func NewTokenGenerator(config *config.Config) *TokenGenerator {
	return &TokenGenerator{
		secret: config.Security.Secret,
		ttl:    config.Security.TokenTTL,
	}
}

func (g *TokenGenerator) Generate(claims *TokenClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": claims.Request.URL.String(),
		"sub": claims.User.Name(),
		"exp": time.Now().Add(g.ttl).Unix(),
	})

	return t.SignedString([]byte(g.secret))
}

type TokenExtractor struct {
	secret string
}

func NewTokenExtractor(config *config.Config) *TokenExtractor {
	return &TokenExtractor{
		secret: config.Security.Secret,
	}
}

func (e *TokenExtractor) Extract(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(e.secret), nil
	})
}
