package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

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
func (g *Guard) ProtectWithJWT(handler fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		if token, exists := headers["Token"]; !exists {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		} else {
			if err := g.authService.ValidateToken(ctx.Context(), token); err != nil {
				return ctx.SendStatus(fiber.StatusUnauthorized)
			}
		}
		return handler(ctx)
	}
}
