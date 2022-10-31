package wallet

import "github.com/gofiber/fiber/v2"

type WalletService interface {
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
func (h *Handler) GetUserWallets(ctx *fiber.Ctx) error {
	panic("implement me")
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/wallets", h.authGuard.ProtectWithJWT(h.GetUserWallets))
}
