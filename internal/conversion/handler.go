package conversion

import "github.com/gofiber/fiber/v2"

type ConversionService interface {
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
	panic("implement me")
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/conversion")
	appGroup.Post("/offer", h.authGuard.ProtectWithJWT(h.CurrencyConversionOffer))
}
