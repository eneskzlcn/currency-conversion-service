package wallet

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/common"
	"github.com/gofiber/fiber/v2"
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
}

func NewHandler(service WalletService, guard AuthGuard) *Handler {
	if service == nil || guard == nil {
		return nil
	}
	return &Handler{
		walletService: service,
		authGuard:     guard,
	}
}
func (h *Handler) GetUserWalletAccounts(ctx *fiber.Ctx) error {
	userID := ctx.Locals(common.USER_ID_CTX_KEY).(int)
	if userID < 0 {
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
