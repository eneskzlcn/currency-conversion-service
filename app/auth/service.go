package auth

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/app/model"
	"github.com/eneskzlcn/currency-conversion-service/config"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"time"
)

type Repository interface {
	GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (model.User, error)
}
type service struct {
	config     config.Jwt
	repository Repository
	logger     *zap.SugaredLogger
}

func NewService(config config.Jwt, repository Repository, logger *zap.SugaredLogger) *service {
	if repository == nil {
		return nil
	}
	return &service{config: config, repository: repository, logger: logger}
}

func (s *service) Tokenize(ctx context.Context, credentials LoginRequest) (TokenResponse, error) {
	user, err := s.repository.GetUserByUsernameAndPassword(ctx, credentials.Username, credentials.Password)
	if err != nil {
		return TokenResponse{}, err
	}

	tokenDuration := time.Duration(s.config.ATExpirationMinutes) * time.Minute
	expirationTime := time.Now().Add(tokenDuration)

	claims := JWTClaim{
		Username: user.Username,
		UserID:   user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.ATPrivateKey))

	return TokenResponse{
		AccessToken: tokenString,
	}, err
}
func (s *service) ValidateToken(_ context.Context, tokenString string) error {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.config.ATPrivateKey), nil
		},
	)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return errors.New("token expired")
	}
	return nil
}

func (s *service) ExtractUserIDFromToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.config.ATPrivateKey), nil
		},
	)
	if err != nil {
		return -1, err
	}
	if claims, ok := token.Claims.(*JWTClaim); ok {
		return claims.UserID, nil
	}

	return -1, nil
}
