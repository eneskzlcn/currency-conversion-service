package conversion

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ConversionService interface {
	CreateCurrencyConversion(ctx context.Context, userID int, request CurrencyConversionOfferRequest) (bool, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}

type Handler struct {
	conversionService ConversionService
	authGuard         AuthGuard
}

func NewHandler(service ConversionService, guard AuthGuard) *Handler {
	if service == nil || guard == nil {
		return nil
	}
	return &Handler{conversionService: service, authGuard: guard}
}
func (h *Handler) CurrencyConversionOffer(ctx *fiber.Ctx) error {
	var offerRequest CurrencyConversionOfferRequest
	if err := ctx.BodyParser(&offerRequest); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	userID, err := strconv.ParseInt(ctx.Get("userID", "-1"), 10, 32)
	if err != nil || userID < 0 {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	success, err := h.conversionService.
		CreateCurrencyConversion(ctx.Context(), int(userID), offerRequest)

	if err != nil || !success {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/conversion")
	appGroup.Post("/offer", h.authGuard.ProtectWithJWT(h.CurrencyConversionOffer))
}
