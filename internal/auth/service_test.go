package auth_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/auth"
	"github.com/eneskzlcn/currency-conversion-service/internal/config"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/auth"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_NewService(t *testing.T) {
	t.Run("test given config and user repository then it should return new Service when NewService called", func(t *testing.T) {
		givenConfig := config.Jwt{
			ATPrivateKey:        "private",
			ATExpirationSeconds: 200,
		}
		ctrl := gomock.NewController(t)
		mockUserRepository := mocks.NewMockUserRepository(ctrl)
		authService := auth.NewService(givenConfig, mockUserRepository)
		assert.NotNil(t, authService)
	})
	t.Run("test given config and nil user repository then it should not return new Service when NewService called", func(t *testing.T) {
		givenConfig := config.Jwt{
			ATPrivateKey:        "private",
			ATExpirationSeconds: 200,
		}
		authService := auth.NewService(givenConfig, nil)
		assert.Nil(t, authService)
	})
}
func Test_Tokenize(t *testing.T) {
	givenConfig := config.Jwt{
		ATPrivateKey:        "private",
		ATExpirationSeconds: 200,
	}
	ctrl := gomock.NewController(t)
	mockUserRepository := mocks.NewMockUserRepository(ctrl)
	authService := auth.NewService(givenConfig, mockUserRepository)

	t.Run("given existing user credentials then it should return access token when Tokenize called", func(t *testing.T) {
		givenCredentials := auth.UserTokenCredentials{
			Username: "iamexistinguser",
			Password: "iamexistingpassword",
		}
		expectedUser := entity.User{
			ID:        1,
			Username:  "iamexistinguser",
			Password:  "iamexistingpassword",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		ctx := context.Background()
		mockUserRepository.EXPECT().
			GetUserByUsernameAndPassword(ctx, givenCredentials.Username, givenCredentials.Password).
			Return(expectedUser, nil)

		accessToken, err := authService.Tokenize(ctx, givenCredentials)
		assert.Nil(t, err)
		assert.NotEmpty(t, accessToken)
	})
	t.Run("given not existing user credentials then it should return empty access token with error when Tokenize called", func(t *testing.T) {
		givenCredentials := auth.UserTokenCredentials{
			Username: "iamnotexistinguser",
			Password: "iamnotexistingpassword",
		}
		ctx := context.Background()
		mockUserRepository.EXPECT().
			GetUserByUsernameAndPassword(ctx, givenCredentials.Username, givenCredentials.Password).
			Return(entity.User{}, errors.New("user not found"))
		accessToken, err := authService.Tokenize(ctx, givenCredentials)
		assert.NotNil(t, err)
		assert.Empty(t, accessToken)
	})
}
