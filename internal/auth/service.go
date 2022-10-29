package auth

import (
	"context"
	"errors"
	"github.com/eneskzlcn/currency-conversion-service/internal/config"
	"github.com/eneskzlcn/currency-conversion-service/internal/entity"
	"github.com/golang-jwt/jwt"
	"time"
)

type UserRepository interface {
	GetUserByUsernameAndPassword(ctx context.Context, username string, password string) (entity.User, error)
}
type Service struct {
	config         config.Jwt
	userRepository UserRepository
}

func NewService(config config.Jwt, repository UserRepository) *Service {
	if repository == nil {
		return nil
	}
	return &Service{config: config, userRepository: repository}
}

func (s *Service) Tokenize(ctx context.Context, credentials UserTokenCredentials) (TokenResponse, error) {
	user, err := s.userRepository.GetUserByUsernameAndPassword(ctx, credentials.Username, credentials.Password)
	if err != nil {
		return TokenResponse{}, err
	}

	tokenDuration := time.Duration(s.config.ATExpirationSeconds) * time.Second
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
func (s *Service) ValidateToken(ctx context.Context, tokenString string) error {
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
