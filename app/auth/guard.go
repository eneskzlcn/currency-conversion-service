package auth

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Service interface {
	ValidateToken(ctx context.Context, tokenString string) error
	ExtractUserIDFromToken(tokenString string) (int, error)
	Tokenize(ctx context.Context, request LoginRequest) (TokenResponse, error)
}

type guard struct {
	service Service
	logger  *zap.SugaredLogger
}

func NewGuard(service Service, logger *zap.SugaredLogger) *guard {
	return &guard{service: service, logger: logger}
}
func (g *guard) ProtectWithJWT(handler fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headers := ctx.GetReqHeaders()
		token, exists := headers["Token"]
		if !exists {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		if err := g.service.ValidateToken(ctx.Context(), token); err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		userID, _ := g.service.ExtractUserIDFromToken(token)
		ctx.Locals(common.USER_ID_CTX_KEY, userID)
		return handler(ctx)
	}
}
