package exchange

import (
	"context"
	"github.com/gofiber/fiber/v2"
)

type ExchangeService interface {
	CreateExchangeRate(ctx context.Context, request ExchangeRateRequest) (ExchangeRateResponse, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}
type Handler struct {
	exchangeService ExchangeService
	authGuard       AuthGuard
}

func NewHandler(service ExchangeService, guard AuthGuard) *Handler {
	if service == nil || guard == nil {
		return nil
	}
	return &Handler{exchangeService: service, authGuard: guard}
}
func (h *Handler) GetExchangeRate(ctx *fiber.Ctx) error {
	var request ExchangeRateRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	exchangeRate, err := h.exchangeService.CreateExchangeRate(ctx.Context(), request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(fiber.StatusCreated).JSON(exchangeRate)
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/exchange")
	appGroup.Get("/rate", h.authGuard.ProtectWithJWT(h.GetExchangeRate))
}
