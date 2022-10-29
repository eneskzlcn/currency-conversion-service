package auth

import "context"

type AuthService interface {
	ValidateToken(ctx context.Context, tokenString string) error
}
type Guard struct {
	authService AuthService
}

func NewGuard(service AuthService) *Guard {
	if service == nil {
		return nil
	}
	return &Guard{authService: service}
}
