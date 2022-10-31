package wallet

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/app/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WalletService interface {
	GetUserWalletAccounts(ctx context.Context, userID int) (UserWalletAccountsResponse, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}
type Handler struct {
	walletService WalletService
	authGuard     AuthGuard
	logger        *zap.SugaredLogger
}

func NewHandler(service WalletService, guard AuthGuard, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		walletService: service,
		authGuard:     guard,
		logger:        logger,
	}
}
func (h *Handler) GetUserWalletAccounts(ctx *fiber.Ctx) error {
	userID, exists := ctx.Locals(common.USER_ID_CTX_KEY).(int)
	if !exists {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	userWalletAccounts, err := h.walletService.GetUserWalletAccounts(ctx.Context(), userID)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(fiber.StatusOK).JSON(userWalletAccounts)
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/wallets", h.authGuard.ProtectWithJWT(h.GetUserWalletAccounts))
}
