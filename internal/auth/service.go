package auth

import "github.com/eneskzlcn/currency-conversion-service/internal/config"

type Service struct {
	config config.Jwt
}

func NewService(config config.Jwt) *Service {
	return &Service{config: config}
}
func (s *Service) Tokenize(credentials UserTokenCredentials) (TokenResponse, error) {
	panic("implement me")
}
