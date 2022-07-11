package security

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/golang-jwt/jwt/v4"
)

type TokenManager struct {
	secret string
}

func NewTokenManager(config *config.Config) *TokenManager {
	return &TokenManager{
		secret: config.Security.Secret,
	}
}

func (m *TokenManager) Sign(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.secret))
}

func (m *TokenManager) Extract(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.secret), nil
	})
}
