package security_test

import (
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestNewTokenGenerator(t *testing.T) {
	generator := security.NewTokenGenerator(&config.Config{
		Security: config.SecurityConfig{
			Secret:   "test",
			TokenTTL: time.Minute,
		},
	})

	request, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	claims := &security.TokenClaims{
		Request: request,
		User:    security.NewUser("test", nil, []security.Scope{}),
	}

	token, err := generator.Generate(claims)
	if err != nil {
		t.Fatal(err)
	}

	require.NotNil(t, token)
}

func TestNewTokenExtractor(t *testing.T) {
	extractor := security.NewTokenExtractor(&config.Config{
		Security: config.SecurityConfig{
			Secret: "test",
		},
	})

	t.Run("failure", func(t *testing.T) {
		token, err := extractor.Extract("")

		require.Error(t, err)
		require.Nil(t, token)
	})

	t.Run("success", func(t *testing.T) {
		token, err := extractor.Extract("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjk5OTkvYXV0aC90b2tlbiIsInN1YiI6InRlc3QifQ.mij_yuFqRYsFwPrWvQTri0-Lh1WC2_lqlg1ys8h_P38")

		require.NoError(t, err)
		require.NotNil(t, token)
	})
}
