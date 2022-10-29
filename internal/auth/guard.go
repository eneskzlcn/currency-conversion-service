package auth

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type AuthService interface {
	ValidateToken(ctx context.Context, tokenString string) error
	ExtractUserIDFromToken(tokenString string) (int, error)
	Tokenize(ctx context.Context, credentials UserAuthRequest) (TokenResponse, error)
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
		token, exists := headers["Token"]
		if !exists {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		if err := g.authService.ValidateToken(ctx.Context(), token); err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		userID, _ := g.authService.ExtractUserIDFromToken(token)
		ctx.Set("userID", strconv.FormatInt(int64(userID), 10))
		return handler(ctx)
	}
}
