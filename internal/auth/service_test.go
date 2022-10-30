package auth_test

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/auth"
	"github.com/eneskzlcn/currency-conversion-service/internal/config"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	mocks "github.com/eneskzlcn/currency-conversion-service/internal/mocks/auth"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func NewAuthServiceAndMockRepoWithDefaultConfig(t *testing.T) (*auth.Service, *mocks.MockAuthRepository) {
	givenConfig := config.Jwt{
		ATPrivateKey:        "private",
		ATExpirationSeconds: 200,
	}
	ctrl := gomock.NewController(t)
	mockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	return auth.NewService(givenConfig, mockAuthRepository), mockAuthRepository
}
func NewAuthServiceAndMockRepoWithGivenConfig(t *testing.T, config config.Jwt) (*auth.Service, *mocks.MockAuthRepository) {
	ctrl := gomock.NewController(t)
	MockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	return auth.NewService(config, MockAuthRepository), MockAuthRepository
}
func CreateMockValidToken(config config.Jwt, user entity.User) (string, error) {
	tokenDuration := time.Duration(config.ATExpirationSeconds) * time.Second
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
func Test_NewService(t *testing.T) {
	t.Run("test given config and user repository then it should return new Service when NewService called", func(t *testing.T) {
		authService, _ := NewAuthServiceAndMockRepoWithDefaultConfig(t)
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
	authService, mockAuthRepository := NewAuthServiceAndMockRepoWithDefaultConfig(t)

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
		ATExpirationSeconds: 1,
	}
	ctrl := gomock.NewController(t)
	mockAuthRepository := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewService(config, mockAuthRepository)

	t.Run("given valid signed token then it should return nil when ValidateToken called", func(t *testing.T) {
		givenUser := entity.User{
			ID:       1,
			Username: "ex",
			Password: "ex",
		}
		token, err := CreateMockValidToken(config, givenUser)
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
		token, err := CreateMockValidToken(config, givenUser)
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
		ATExpirationSeconds: 20,
	}
	authService, _ := NewAuthServiceAndMockRepoWithGivenConfig(t, givenConfig)
	t.Run("given valid token then it should extract the user id from token when ExtractUserIDFromToken called", func(t *testing.T) {
		givenUser := entity.User{
			ID:       2,
			Username: "asd",
		}
		token, err := CreateMockValidToken(givenConfig, givenUser)
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
