package auth_test

import (
	"github.com/eneskzlcn/currency-conversion-service/internal/auth"
	"github.com/eneskzlcn/currency-conversion-service/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGivenJWTConfigThenItShouldReturnNewServiceWhenNewServiceCalled(t *testing.T) {
	givenConfig := config.Jwt{
		ATPrivateKey:        "private",
		ATExpirationSeconds: 200,
	}
	authService := auth.NewService(givenConfig)
	assert.NotNil(t, authService)
}

func Test_Tokenize(t *testing.T) {
	givenConfig := config.Jwt{
		ATPrivateKey:        "private",
		ATExpirationSeconds: 200,
	}
	authService := auth.NewService(givenConfig)

	t.Run("given existing user credentials then it should return access token when Tokenize called", func(t *testing.T) {
		givenCredentials := auth.UserTokenCredentials{
			Username: "iamexistinguser",
			Password: "iamexistingpassword",
		}
		accessToken, err := authService.Tokenize(givenCredentials)
		assert.Nil(t, err)
		assert.NotEmpty(t, accessToken)
	})
	t.Run("given not existing user credentials then it should return empty access token with error when Tokenize called", func(t *testing.T) {
		givenCredentials := auth.UserTokenCredentials{
			Username: "iamnotexistinguser",
			Password: "iamnotexistingpassword",
		}
		accessToken, err := authService.Tokenize(givenCredentials)
		assert.NotNil(t, err)
		assert.Empty(t, accessToken)
	})
}
