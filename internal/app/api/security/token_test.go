package security_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestTokenManager_Sign(t *testing.T) {
	manager := security.NewTokenManager(&config.Config{
		Security: config.SecurityConfig{
			Secret: "test",
		},
	})

	request, err := http.NewRequest(http.MethodGet, "/", nil)

	require.NoError(t, err)
	require.NotNil(t, request)

	claims := jwt.MapClaims{
		"iss": request.URL.String(),
		"sub": "test",
		"exp": time.Now().Add(5 * time.Second).Unix(),
	}

	accessToken, err := manager.Sign(claims)

	require.NoError(t, err)
	require.NotNil(t, accessToken)
}

func TestTokenManager_Extract(t *testing.T) {
	manager := security.NewTokenManager(&config.Config{
		Security: config.SecurityConfig{
			Secret: "test",
		},
	})

	t.Run("failure", func(t *testing.T) {
		token, err := manager.Extract("")

		require.Error(t, err)
		require.Nil(t, token)
	})

	t.Run("success", func(t *testing.T) {
		token, err := manager.Extract("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjk5OTkvYXV0aC90b2tlbiIsInN1YiI6InRlc3QifQ.mij_yuFqRYsFwPrWvQTri0-Lh1WC2_lqlg1ys8h_P38")

		require.NoError(t, err)
		require.NotNil(t, token)
	})
}
