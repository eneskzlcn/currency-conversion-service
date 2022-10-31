package auth

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthService interface {
	ValidateToken(ctx context.Context, tokenString string) error
	ExtractUserIDFromToken(tokenString string) (int, error)
	Tokenize(ctx context.Context, request LoginRequest) (TokenResponse, error)
}

type Guard struct {
	authService AuthService
	logger      *zap.SugaredLogger
}

func NewGuard(service AuthService, logger *zap.SugaredLogger) *Guard {
	return &Guard{authService: service, logger: logger}
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
		g.logger.Debugf("Authorized user with ID:%d", userID)
		ctx.Locals(common.USER_ID_CTX_KEY, userID)
		return handler(ctx)
	}
}
