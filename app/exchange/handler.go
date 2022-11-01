package exchange

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ExchangeService interface {
	PrepareExchangeRateOffer(ctx context.Context, request ExchangeRateRequest) (ExchangeRateResponse, error)
}
type AuthGuard interface {
	ProtectWithJWT(handler fiber.Handler) fiber.Handler
}
type Handler struct {
	exchangeService ExchangeService
	authGuard       AuthGuard
	logger          *zap.SugaredLogger
}

func NewHandler(service ExchangeService, guard AuthGuard, logger *zap.SugaredLogger) *Handler {
	return &Handler{exchangeService: service, authGuard: guard, logger: logger}
}

//GetExchangeRate godoc
//@Summary Create an exchange rate offer
//@Description creates an exchange rate offer for given currencies
//@Param exchangeRateRequest body ExchangeRateRequest true "body params"
//@Param accessToken header string true "header params"
//@Accept  json
//@Produce  json
//@Success 200 {object} ExchangeRateResponse
//@Failure 400
//@Failure 401 {string} string "Unauthorized"
//@Failure 404
//@Failure 500
//@Router /exchange/rate [get]
func (h *Handler) GetExchangeRate(ctx *fiber.Ctx) error {
	var request ExchangeRateRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	exchangeRate, err := h.exchangeService.PrepareExchangeRateOffer(ctx.Context(), request)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(fiber.StatusCreated).JSON(exchangeRate)
}
func (h *Handler) RegisterRoutes(app *fiber.App) {
	appGroup := app.Group("/exchange")
	appGroup.Get("/rate", h.authGuard.ProtectWithJWT(h.GetExchangeRate))
}
