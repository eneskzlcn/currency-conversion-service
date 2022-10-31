package conversion

import (
	"context"
	"github.com/eneskzlcn/currency-conversion-service/internal/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ConversionService interface {
	ConvertCurrencies(ctx context.Context, userID int, request CurrencyConversionOfferRequest) (bool, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}

type Handler struct {
	conversionService ConversionService
	authGuard         AuthGuard
	logger            *zap.SugaredLogger
}

func NewHandler(service ConversionService, guard AuthGuard, logger *zap.SugaredLogger) *Handler {
	return &Handler{conversionService: service, authGuard: guard, logger: logger}
}
func (h *Handler) CurrencyConversionOffer(ctx *fiber.Ctx) error {
	userID, exists := ctx.Locals(common.USER_ID_CTX_KEY).(int)
	h.logger.Infof("Currency Conversion Offer Request Arrived. User ID: %d", userID)
	if !exists {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	var offerRequest CurrencyConversionOfferRequest
	if err := ctx.BodyParser(&offerRequest); err != nil {
		h.logger.Debug("Can not parse request body")
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	success, err := h.conversionService.
		ConvertCurrencies(ctx.Context(), userID, offerRequest)

	if err != nil || !success {
		h.logger.Debug("Convert currency operation ended up with error")
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.SendStatus(fiber.StatusOK)
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/conversion")
	appGroup.Post("/offer", h.authGuard.ProtectWithJWT(h.CurrencyConversionOffer))
}
