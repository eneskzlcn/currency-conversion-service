//go:build unit

package auth_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/auth"
	"github.com/eneskzlcn/currency-conversion-service/app/entity"
	mocks "github.com/eneskzlcn/currency-conversion-service/app/mocks/auth"
	"github.com/eneskzlcn/currency-conversion-service/config"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func Test_Tokenize(t *testing.T) {
	authService, mockAuthRepository := newAuthServiceAndMockRepoWithDefaultConfig(t)

	t.Run("given existing user credentials then it should return access token when Tokenize called", func(t *testing.T) {
		givenCredentials := auth.LoginRequest{
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
		mockAuthRepository.EXPECT().
			GetUserByUsernameAndPassword(ctx, givenCredentials.Username, givenCredentials.Password).
			Return(expectedUser, nil)

		accessToken, err := authService.Tokenize(ctx, givenCredentials)
		assert.Nil(t, err)
		assert.NotEmpty(t, accessToken)
	})
	t.Run("given not existing user credentials then it should return empty access token with error when Tokenize called", func(t *testing.T) {
		givenCredentials := auth.LoginRequest{
			Username: "iamnotexistinguser",
			Password: "iamnotexistingpassword",
		}
		ctx := context.Background()
		mockAuthRepository.EXPECT().
			GetUserByUsernameAndPassword(ctx, givenCredentials.Username, givenCredentials.Password).
			Return(entity.User{}, errors.New("user not found"))
		accessToken, err := authService.Tokenize(ctx, givenCredentials)
		assert.NotNil(t, err)
		assert.Empty(t, accessToken)
	})
}

func Test_ValidateToken(t *testing.T) {
	config := config.Jwt{
		ATPrivateKey:        "privateKey",
		ATExpirationMinutes: 1,
	}
	ctrl := gomock.NewController(t)
	mockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewService(config, mockAuthRepository, zap.S())

	t.Run("given valid signed token then it should return nil when ValidateToken called", func(t *testing.T) {
		givenUser := entity.User{
			ID:       1,
			Username: "ex",
			Password: "ex",
		}
		token, err := createMockValidToken(config, givenUser)
		assert.Nil(t, err)
		assert.NotEmpty(t, token)

		err = authService.ValidateToken(context.Background(), token)
		assert.Nil(t, err)

	})
	t.Run("given not valid token then it should return error when ValidateToken called", func(t *testing.T) {
		token := ""
		err := authService.ValidateToken(context.Background(), token)
		assert.NotNil(t, err)
	})
	t.Run("given valid but expired token then it should return error when ValidateToken called", func(t *testing.T) {
		givenUser := entity.User{
			ID:       1,
			Username: "ex",
			Password: "ex",
		}
		token, err := createMockValidToken(config, givenUser)
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		time.Sleep(2 * time.Second)
		err = authService.ValidateToken(context.Background(), token)
		assert.NotNil(t, err)
	})
}

func TestService_ExtractUserIDFromToken(t *testing.T) {
	givenConfig := config.Jwt{
		ATPrivateKey:        "mypr",
		ATExpirationMinutes: 20,
	}
	authService, _ := newAuthServiceAndMockRepoWithGivenConfig(t, givenConfig)
	t.Run("given valid token then it should extract the user id from token when ExtractUserIDFromToken called", func(t *testing.T) {
		givenUser := entity.User{
			ID:       2,
			Username: "asd",
		}
		token, err := createMockValidToken(givenConfig, givenUser)
		assert.Nil(t, err)

		userID, err := authService.ExtractUserIDFromToken(token)
		assert.Nil(t, err)
		assert.Equal(t, 2, userID)
	})
	t.Run("given invalid token then it should return -1 with error when ExtractUserIDFromToken called", func(t *testing.T) {
		userID, err := authService.ExtractUserIDFromToken("not valid token")
		assert.NotNil(t, err)
		assert.Equal(t, -1, userID)
	})
}

func newAuthServiceAndMockRepoWithDefaultConfig(t *testing.T) (*auth.Service, *mocks.MockAuthRepository) {
	givenConfig := config.Jwt{
		ATPrivateKey:        "private",
		ATExpirationMinutes: 200,
	}
	ctrl := gomock.NewController(t)
	mockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	return auth.NewService(givenConfig, mockAuthRepository, zap.S()), mockAuthRepository
}
func newAuthServiceAndMockRepoWithGivenConfig(t *testing.T, config config.Jwt) (*auth.Service, *mocks.MockAuthRepository) {
	ctrl := gomock.NewController(t)
	MockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	return auth.NewService(config, MockAuthRepository, zap.S()), MockAuthRepository
}
func createMockValidToken(config config.Jwt, user entity.User) (string, error) {
	tokenDuration := time.Duration(config.ATExpirationMinutes) * time.Second
	expirationTime := time.Now().Add(tokenDuration)
	claims := auth.JWTClaim{
		Username: user.Username,
		UserID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.ATPrivateKey))
}
