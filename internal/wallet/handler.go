package wallet

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
	userID, err := strconv.ParseInt(ctx.Get("userID", "-1"), 10, 32)
	if err != nil || userID < 0 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	userWalletAccounts, err := h.walletService.GetUserWalletAccounts(ctx.Context(), int(userID))
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(fiber.StatusOK).JSON(userWalletAccounts)
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/wallets", h.authGuard.ProtectWithJWT(h.GetUserWalletAccounts))
}
