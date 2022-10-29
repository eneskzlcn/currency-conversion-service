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
